package log

import (
	"net/http"

	"time"
)

type FormatterRegistry struct {
	ConsoleFormatter Formatter
	FileFormmater    Formatter
	RemoteFormatter  Formatter
}
type RemoteConfig struct {
	EndPoint  string
	Method    string
	Header    http.Header
	Transport *http.RoundTripper
}

type Config struct {
	Location          *time.Location
	Level             LogLevel
	OutputMode        OutputMode
	EntrySize         int
	FileConfig        *FileConfig
	RemoteConfig      *RemoteConfig
	StandardFormatter Formatter
	FormatterRegistry *FormatterRegistry
}

type FileCreateMode int

const (
	DAILYMODE FileCreateMode = iota
	HOURLYMODE
)

type FileConfig struct {
	FileName       string
	LogPath        string
	FileCreateMode FileCreateMode
}

func convertOptions(config *Config) []LogOption {
	var opts []LogOption
	if config.Level != 0 {
		opts = append(opts, WithLevel(config.Level))
	}
	if config.EntrySize != 0 {
		opts = append(opts, WithEntrySize(config.EntrySize))
	}
	if config.OutputMode != 0 {
		opts = append(opts, withOutputMode(config.OutputMode))
	}
	if config.FileConfig != nil {
		opts = append(opts, WithFileMode(config.FileConfig.FileName, config.FileConfig.LogPath, config.FileConfig.FileCreateMode))
	}
	if config.RemoteConfig != nil {
		opts = append(opts, WithRemoteMode(config.RemoteConfig.EndPoint, config.RemoteConfig.Method, config.RemoteConfig.Header, config.RemoteConfig.Transport))
	}
	if config.Location != nil {
		opts = append(opts, WithLocation(config.Location))
	}
	if config.StandardFormatter != nil {
		opts = append(opts, WithStandardFormatter(config.StandardFormatter))
	}
	if config.FormatterRegistry.ConsoleFormatter != nil {
		opts = append(opts, WithConsoleFormatter(config.FormatterRegistry.ConsoleFormatter))
	}
	if config.FormatterRegistry.FileFormmater != nil {
		opts = append(opts, WithFileModeFormatter(config.FormatterRegistry.FileFormmater))
	}
	if config.FormatterRegistry.RemoteFormatter != nil {
		opts = append(opts, WithRemoteFormatter(config.FormatterRegistry.RemoteFormatter))
	}
	return opts
}
