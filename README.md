# go-log

go-log is a lightweight and flexible logging library for Go, designed to support multiple output modes and customizable formatting.

## Features
- **Multiple Output Modes**: Console, file, and remote logging.
- **Customizable Formatters**: Define your own log entry format.
- **Log Levels**: DEBUG, INFO, WARN, ERROR.
- **Thread-Safe**: Designed for concurrent use.
- **Dynamic Configuration**: Adjust log levels and output modes at runtime.

## Installation
```bash
go get github.com/winey-dev/go-log
```

## Quick Start
### Default Logger
```go
package main

import (
	"github.com/winey-dev/go-log"
)

func main() {
	// Create a new logger with default settings (INFO level, console output)
	mlog, err := log.NewLogger("my-app")
	if err != nil {
		panic(err)
	}
	defer mlog.Close()

	mlog.Info("Hello, World!")
	mlog.Debug("This is a debug message")
	mlog.Warn("This is a warning")
	mlog.Error("This is an error")
}
```

### Change Log Level
```go
mlog.SetLogLevel(log.DEBUG)
```

### Console and File Logging
```go
mlog, err := log.NewLogger("my-app",
	log.WithFileMode("", "", log.DAILYMODE), // Enable file logging
)
if err != nil {
	panic(err)
}
defer mlog.Close()
```

### File-Only Logging
```go
mlog, err := log.NewLogger("my-app",
	log.WithFileMode("", "", log.DAILYMODE),
	log.WithConsoleModeOff(), // Disable console logging
)
if err != nil {
	panic(err)
}
defer mlog.Close()
```

### Custom Formatter
```go
mlog, err := log.NewLogger("my-app",
	log.WithConsoleFormatter(func(t time.Time, level log.LogLevel, format string, args ...any) string {
		return fmt.Sprintf("%s [%s] %s", t.Format("2006-01-02 15:04:05"), log.LoglevelNames[level], fmt.Sprintf(format, args...))
	}),
)
if err != nil {
	panic(err)
}
defer mlog.Close()
```

## License
This project is licensed under the Apache 2.0 License. See the LICENSE file for details.