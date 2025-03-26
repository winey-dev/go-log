package log

import (
	"io"
	"net/http"
	"time"
)

// LogOption is a function that sets the logger configuration.
type LogOption func(*logger)

// globalLogger is used in the option. It is set when the WithGlobal option is executed.
func WithGlobal() LogOption {
	return func(l *logger) {
		globalLogger = l
	}
}

// WithLocation sets the location of the logger. The default is time.Local.
func WithLocation(location *time.Location) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.Location = location
	}
}

// WithLevel sets the log level of the logger. The default is INFO.
func WithLevel(level LogLevel) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.Level = level
	}
}

// WithEntrySize sets the size of the log entry. The default is 4096.
func WithEntrySize(size int) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.EntrySize = size
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

func WithConsoleOutPut(w io.Writer) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		if l.config.ConsoleConfig == nil {
			l.config.ConsoleConfig = &ConsoleConfig{}
		}
		l.config.ConsoleConfig.Writer = w
	}
}

// WithFileMode sets the file name, log path, and file create mode of the logger.
// The default is an empty string for the file name and log path and DAILYMODE for the file create mode.
// The file name is the name of the log file.
// The log path is the path to the log file.
// The file create mode is the mode of the log file.
// The file create mode can be DAILYMODE or HOURLYMODE.
// The file create mode is used to create a new log file every day or every hour.
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

// WithRemoteMode sets the endpoint, method, header, and transport of the logger.
// The default is an empty string for the endpoint and method, nil for the header, and nil for the transport.
// The endpoint is the endpoint of the remote logger.
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

// WithConsoleModeOff sets the console mode of the logger to off.
func WithConsoleModeOff() LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.OutputMode &= ^OutputModeConsole

	}
}

// WithFileModeOff sets the file mode of the logger to off.
func WithStandardFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.StandardFormatter = formatter
	}
}

// WithConsoleFormatter sets the console formatter of the logger.
// The default is the standard formatter.
func WithConsoleFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.ConsoleFormatter = formatter
	}
}

// WithFileModeFormatter sets the file mode formatter of the logger.
// The default is the standard formatter.
func WithFileModeFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.FileFormmater = formatter
	}
}

// WithRemoteFormatter sets the remote formatter of the logger.
// The default is the standard formatter.
func WithRemoteFormatter(formatter Formatter) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry.RemoteFormatter = formatter
	}
}

// WithRegisterFormatter sets the formatter register of the logger.
// The default is the standard formatter for each mode.
// This function allows setting all formatters at once instead of calling WithConsoleFormatter,
// WithFileModeFormatter, and WithRemoteFormatter separately.
func WithRegisterFormatter(formatterRegister *FormatterRegistry) LogOption {
	return func(l *logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.config.FormatterRegistry = formatterRegister
	}
}
