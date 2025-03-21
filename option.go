package log

import (
	"net/http"
	"time"
)

type LogOption func(*logger)

// globalLogger is used in the option. It is set when the WithGlobal option is executed.

func WithGlobal() LogOption {
	return func(l *logger) {
		globalLogger = l
	}
}

func WithLevel(level LogLevel) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.Level = level
	}
}

func withOutputMode(mode OutputMode) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.OutputMode = mode
	}
}

func WithConsoleMode() LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.OutputMode |= OutputModeConsole
	}
}

func WithFileMode(fileName, logPath string, mode FileCreateMode) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		if l.config.FileConfig == nil {
			l.config.FileConfig = &FileConfig{}
		}
		l.config.FileConfig.FileName = fileName
		l.config.FileConfig.LogPath = logPath
		l.config.FileConfig.FileCreateMode = mode
		l.config.OutputMode |= OutputModeFile
	}
}

func WithRemoteMode(endpoint, method string, header http.Header, transport *http.RoundTripper) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		if l.config.RemoteConfig == nil {
			l.config.RemoteConfig = &RemoteConfig{}
		}
		l.config.RemoteConfig.Method = method
		l.config.RemoteConfig.EndPoint = endpoint
		l.config.RemoteConfig.Header = header
		l.config.RemoteConfig.Transport = transport
		l.config.OutputMode |= OutputModeRemote
	}
}

func WithConsoleModeOff() LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.OutputMode &= ^OutputModeConsole

	}
}

func WithLocation(location *time.Location) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.Location = location
	}
}

func WithStandardFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.StandardFormatter = formatter
	}
}

func WithConsoleFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.ConsoleFormatter = formatter
	}
}

func WithFileModeFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.FileModeFormmater = formatter
	}
}

func WithRemoteFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.RemoteFormatter = formatter
	}
}

func WithRegisterFormatter(formatterRegister *FormatterRegistry) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry = formatterRegister
	}
}
