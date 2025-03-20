package log

import (
	"sync"
	"time"
)

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

type logger struct {
	name          string
	config        *Config
	dynamicWriter *dynamicWriter
	mtx           sync.RWMutex
}

var globalLogger *logger

type LogLevel int

const (
	NONE  LogLevel = 0
	DEBUG LogLevel = 1 + iota
	INFO           // 2
	WARN           // 3
	ERROR          // 4
)

var loglevelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}
var _ = loglevelNames

type OutputMode uint8

const (
	OutputModeConsole OutputMode = 1 << iota
	OutputModeFile               // 2
	OutputModeRemote             // 4
)

func NewLoggerFormConfig(name string, config *Config) (*logger, error) {
	return NewLogger(name, convertOptions(config)...)
}

func NewLogger(name string, opts ...LogOption) (*logger, error) {
	logger := &logger{
		name: name,
		config: &Config{
			Level:             INFO,
			OutputMode:        OutputModeConsole,
			Location:          time.Local,
			StandardFormatter: defaultFormatter,
			FormatterRegistry: &FormatterRegistry{},
		},
		dynamicWriter: &dynamicWriter{},
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger.init()
}

func (l *logger) Debug(format string, args ...any) {
	l.logf(DEBUG, format, args...)
}
func (l *logger) Info(format string, args ...any) {
	l.logf(INFO, format, args...)
}
func (l *logger) Warn(format string, args ...any) {
	l.logf(WARN, format, args...)
}
func (l *logger) Error(format string, args ...any) {
	l.logf(ERROR, format, args...)
}

func (l *logger) logf(level LogLevel, format string, args ...any) {
	if level < l.config.Level {
		return
	}
	now := time.Now().In(l.config.Location)
	entry := &logEntry{
		t:      now,
		level:  level,
		format: format,
		args:   args,
	}
	l.dynamicWriter.ch <- entry
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.config.Level = level
}

func Debug(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Debug(format, args...)
}

func Info(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(format, args...)
}
func Warn(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Warn(format, args...)
}

func Error(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Error(format, args...)
}

func SetLogLevel(level LogLevel) {
	if globalLogger == nil {
		return
	}
	globalLogger.SetLogLevel(level)
}

func (l *logger) init() (*logger, error) {
	// 각 모드별 formatter 설정이 없을 경우 기본 포매터 사용
	if l.config.OutputMode&OutputModeConsole != 0 {
		if l.config.FormatterRegistry.ConsoleFormatter == nil {
			l.config.FormatterRegistry.ConsoleFormatter = l.config.StandardFormatter
		}
	}

	if l.config.OutputMode&OutputModeFile != 0 {
		if l.config.FormatterRegistry.FileModeFormmater == nil {
			l.config.FormatterRegistry.FileModeFormmater = l.config.StandardFormatter
		}
		if l.config.FileConfig == nil {
			l.config.FileConfig = &FileConfig{
				FileName:       l.name,
				LogPath:        "log", // $HOME/log
				FileCreateMode: DAILYMODE,
			}
		} else {
			if l.config.FileConfig.FileName == "" {
				l.config.FileConfig.FileName = l.name
			}
			if l.config.FileConfig.LogPath == "" {
				l.config.FileConfig.LogPath = "log"
			}
		}
	}

	if l.config.OutputMode&OutputModeRemote != 0 {
		if l.config.FormatterRegistry.RemoteFormatter == nil {
			l.config.FormatterRegistry.RemoteFormatter = l.config.StandardFormatter
		}
		if l.config.RemoteConfig == nil {
			return nil, ErrRemoteConfig
		} else if l.config.RemoteConfig.EndPoint == "" { // remote addr is required
			return nil, ErrRemoteEndpoint
		}
	}

	l.dynamicWriter = newDynamicWriter(l)
	l.dynamicWriter.run()
	return l, nil
}

func (l *logger) Close() {
	l.dynamicWriter.close()
}
