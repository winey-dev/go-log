package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/winey-dev/go-log"
)

func TestLogConsole(t *testing.T) {
	mlog, _ := log.NewLogger("test",
		log.WithLevel(log.DEBUG),
		log.WithConsoleMode(),
	)
	mlog.Debug("Test Debug\n")
	mlog.Info("Test Info\n")
	mlog.Warn("Test Warn\n")
	mlog.Error("Test Error\n")
	mlog.Close()

}

func TestLogConsoleFormatterJSON(t *testing.T) {
	mlog, _ := log.NewLogger("test",
		log.WithLevel(log.DEBUG),
		log.WithConsoleMode(),
		log.WithConsoleFormatter(func(t time.Time, level log.LogLevel, format string, args ...any) string {
			type JSONFORMAT struct {
				Time    string       `json:"time"`
				Level   log.LogLevel `json:"level"`
				Message string       `json:"message"`
			}

			var jsonFormat JSONFORMAT
			jsonFormat.Time = t.Local().String()
			jsonFormat.Level = level
			jsonFormat.Message = fmt.Sprintf(format, args...)
			dat, _ := json.Marshal(&jsonFormat)
			return fmt.Sprintf("%s\n", dat)
		}),
	)
	mlog.Debug("Test Debug\n")
	mlog.Info("Test Info\n")
	mlog.Warn("Test Warn\n")
	mlog.Error("Test Error\n")
	mlog.Close()

}
