package log

import (
	"fmt"
	"time"
)

type Formatter func(t time.Time, level LogLevel, format string, args ...any) string

func defaultFormatter(t time.Time, level LogLevel, format string, args ...any) string {
	return fmt.Sprintf("%s [%5s] %s", t.In(time.Local).Format(time.TimeOnly), loglevelNames[level], fmt.Sprintf(format, args...))
}
