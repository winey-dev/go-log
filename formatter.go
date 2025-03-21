package log

import (
	"fmt"
	"time"
)

// Formatter is a function that formats log entries.
// The Formatter function is called by the logger to format log entries.
type Formatter func(t time.Time, level LogLevel, format string, args ...any) string

func defaultFormatter(t time.Time, level LogLevel, format string, args ...any) string {
	return fmt.Sprintf("%s [%5s] %s", t.In(time.Local).Format(time.TimeOnly+".000"), LoglevelNames[level], fmt.Sprintf(format, args...))
}

// custom formatter is a function that formats log entries.
// example:
// func CustomFormatter(t time.Time, level LogLevel, format string, args ...any) string {
// 	   var customLogEntry struct {
// 		   Time string
// 		   Level string
// 		   Message string
// 	   } {
// 		   Time: t.UTC.String(),
// 		   Level: log.LoglevelNames[level],
// 		   Message: fmt.Sprintf(format, args...),
// 	   }
//     dat, _ := json.Marshal(customLogEntry)
// 	   return fmt.Sprintf("%s\n", string(dat))
// }
//
// ...
// mlog, err := log.NewLogger("my-app",
// 		log.WithConosleFormatter(log.CustomFormatter))
// ...
// mlog.Info("Hello, World!\n")
// mlog.Debug("Hello, World!\n")
// mlog.Warn("Hello, World!\n")
// mlog.Error("Hello, World!\n")
// ...
// output:
// {"Time":"2021-09-01T00:00:00Z","Level":"INFO","Message":"Hello, World!\n"}
// {"Time":"2021-09-01T00:00:00Z","Level":"DEBUG","Message":"Hello, World!\n"}
// {"Time":"2021-09-01T00:00:00Z","Level":"WARN","Message":"Hello, World!\n"}
// {"Time":"2021-09-01T00:00:00Z","Level":"ERROR","Message":"Hello, World!\n"}
