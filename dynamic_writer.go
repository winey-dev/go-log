package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Writer interface {
	Write(t time.Time, level LogLevel, format string, args ...any) (n int, err error)
}

type logEntry struct {
	t      time.Time
	level  LogLevel
	format string
	args   []any
}
type dynamicWriter struct {
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	writers map[OutputMode]Writer
	ch      chan *logEntry
}

func newDynamicWriter(l *logger) *dynamicWriter {
	ctx, cancle := context.WithCancel(context.Background())
	writer := &dynamicWriter{
		ctx:     ctx,
		cancel:  cancle,
		writers: make(map[OutputMode]Writer),
		ch:      make(chan *logEntry, l.config.EntrySize),
	}

	if l.config.OutputMode&OutputModeConsole != 0 {
		writer.writers[OutputModeConsole] = newConsoleWriter(l)
	}

	if l.config.OutputMode&OutputModeFile != 0 {
		writer.writers[OutputModeFile] = newFileWriter(l)
	}

	if l.config.OutputMode&OutputModeRemote != 0 {
		writer.writers[OutputModeRemote] = newRemoteWriter(l)
	}

	return writer
}

func (d *dynamicWriter) run() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			select {
			case logEntry, ok := <-d.ch:
				if !ok {
					return
				}
				d.writer(logEntry.t, logEntry.level, logEntry.format, logEntry.args...)
			case <-d.ctx.Done():
				return
			}
		}
	}()
}

func (d *dynamicWriter) writer(t time.Time, level LogLevel, format string, args ...any) {
	for _, writer := range d.writers {
		_, _ = writer.Write(t, level, format, args...)
	}
}

func (d *dynamicWriter) close() {
	d.cancel()
	d.wg.Wait()
	close(d.ch)

	for entry := range d.ch {
		d.writer(entry.t, entry.level, entry.format, entry.args...)
	}

	for _, writer := range d.writers {
		if f, ok := writer.(*fileWriter); ok {
			_ = f.file.Close()
		}
	}

}

type consoleWriter struct {
	formatter Formatter
}

func newConsoleWriter(l *logger) Writer {
	return &consoleWriter{
		formatter: l.config.FormatterRegistry.ConsoleFormatter,
	}
}

func (c *consoleWriter) Write(t time.Time, level LogLevel, format string, args ...any) (n int, err error) {
	return fmt.Fprint(os.Stdout, c.formatter(t, level, format, args...))
}

type fileWriter struct {
	name            string
	logPath         string
	mode            FileCreateMode
	formatter       Formatter
	currentFileName string
	file            *os.File
}

func newFileWriter(l *logger) Writer {
	// logPath dir를 생성

	var logPath string

	path := os.ExpandEnv(l.config.FileConfig.LogPath)

	if filepath.IsAbs(path) {
		logPath = path
	} else {
		var err error
		logPath, err = filepath.Abs(path)
		if err != nil {
			logPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), path)
		}
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.MkdirAll(logPath, 0755)
	}

	return &fileWriter{
		name:      l.name,
		logPath:   logPath,
		mode:      l.config.FileConfig.FileCreateMode,
		formatter: l.config.FormatterRegistry.FileFormmater,
	}
}

func (f *fileWriter) generatedFileName(t time.Time) string {
	if f.mode == DAILYMODE {
		return fmt.Sprintf("%s/%s.%s.log", f.logPath, f.name, t.Format(time.DateOnly))
	}
	return fmt.Sprintf("%s/%s.%s.log", f.logPath, f.name, t.Format("2006-01-02-15"))
}

func (f *fileWriter) Write(t time.Time, level LogLevel, format string, args ...any) (n int, err error) {
	generatedFileName := f.generatedFileName(t)
	if generatedFileName != f.currentFileName {
		f.currentFileName = generatedFileName
		if f.file != nil {
			_ = f.file.Close()
		}
		f.file, err = os.OpenFile(f.currentFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, err
		}
	}

	if f.file != nil {
		return fmt.Fprint(f.file, f.formatter(t, level, format, args...))
	}

	return 0, nil
}

type remoteWriter struct {
	endpoint  string
	method    string
	header    http.Header
	transport *http.RoundTripper
	formatter Formatter
}

func newRemoteWriter(l *logger) Writer {
	return &remoteWriter{
		method:    l.config.RemoteConfig.Method,
		endpoint:  l.config.RemoteConfig.EndPoint,
		header:    l.config.RemoteConfig.Header,
		transport: l.config.RemoteConfig.Transport,
		formatter: l.config.FormatterRegistry.RemoteFormatter,
	}
}

func (r *remoteWriter) Write(t time.Time, level LogLevel, format string, args ...any) (n int, err error) {
	var client *http.Client
	type remoteLog struct {
		Time    time.Time `json:"time"`
		Level   string    `json:"level"`
		Message string    `json:"message"`
	}

	var log = &remoteLog{
		Time:    t,
		Level:   loglevelNames[level],
		Message: fmt.Sprintf(format, args...),
	}

	go func() {
		dat, err := json.Marshal(log)
		if err != nil {
			return
		}

		buffer := bytes.NewBuffer(dat)

		req, err := http.NewRequest(r.method, r.endpoint, buffer)
		if err != nil {
			return
		}

		if (r.header != nil) && (len(r.header) > 0) {
			req.Header = r.header
		}

		req.Header.Set("Content-Type", "application/json")
		if r.transport != nil {
			client = &http.Client{Transport: *r.transport}
		} else {
			client = http.DefaultClient
		}

		_, err = client.Do(req)
		if err != nil {
			return
		}
	}()
	return 0, nil
}
