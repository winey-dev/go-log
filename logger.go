package log

import (
	"sync"
	"time"
)

// Logger is the interface that wraps the basic logging methods.
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

var LoglevelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

type OutputMode uint8

const (
	// OutputModeConsole is the output mode for the console.
	OutputModeConsole OutputMode = 1 << iota
	// OutputModeFile is the output mode for the file.
	OutputModeFile // 2
	// OutputModeRemote is the output mode for the remote.
	OutputModeRemote // 4
)

// NewLoggerFormConfig creates a new logger from the configuration.
func NewLoggerFormConfig(name string, config *Config) (*logger, error) {
	return NewLogger(name, convertOptions(config)...)
}

// NewLogger creates a new logger with the options.
func NewLogger(name string, opts ...LogOption) (*logger, error) {
	logger := &logger{
		name: name,
		config: &Config{
			Location:          time.Local,
			Level:             INFO,
			OutputMode:        OutputModeConsole,
			EntrySize:         4096,
			StandardFormatter: defaultFormatter,
			FormatterRegistry: &FormatterRegistry{},
		},
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger.init()
}

// Debug logs a message with the DEBUG level.
func (l *logger) Debug(format string, args ...any) {
	l.logf(DEBUG, format, args...)
}

// Info logs a message with the INFO level.
func (l *logger) Info(format string, args ...any) {
	l.logf(INFO, format, args...)
}

// Warn logs a message with the WARN level.
func (l *logger) Warn(format string, args ...any) {
	l.logf(WARN, format, args...)
}

// Error logs a message with the ERROR level.
func (l *logger) Error(format string, args ...any) {
	l.logf(ERROR, format, args...)
}

// logf logs a message with the level and format.
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

// SetLogLevel sets the log level of the logger.
func (l *logger) SetLogLevel(level LogLevel) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.config.Level = level
}

// Debug logs a message with the DEBUG level. It is a wrapper for the global logger.
func Debug(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Debug(format, args...)
}

// Info logs a message with the INFO level. It is a wrapper for the global logger.
func Info(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(format, args...)
}

// Warn logs a message with the WARN level. It is a wrapper for the global logger.
func Warn(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Warn(format, args...)
}

// Error logs a message with the ERROR level. It is a wrapper for the global logger.
func Error(format string, args ...any) {
	if globalLogger == nil {
		return
	}
	globalLogger.Error(format, args...)
}

// SetLogLevel sets the log level of the global logger.
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
		if l.config.FormatterRegistry.FileFormmater == nil {
			l.config.FormatterRegistry.FileFormmater = l.config.StandardFormatter
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

// Close closes the logger.
// It ensures that all remaining log entries in the channel are processed before shutting down.
func (l *logger) Close() {
	l.dynamicWriter.close()
}
