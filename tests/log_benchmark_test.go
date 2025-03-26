package tests

import (
	"testing"

	"github.com/winey-dev/go-log"
)

func BenchmarkLogConsole(b *testing.B) {
	mlog, _ := log.NewLogger("test",
		log.WithLevel(log.DEBUG),
		log.WithConsoleMode(),
		//log.WithConsoleOutPut(io.Discard),
	)
	for i := 0; i < b.N; i++ {
		mlog.Debug("Test Debug\n")
		mlog.Info("Test Info\n")
		mlog.Warn("Test Warn\n")
		mlog.Error("Test Error\n")
	}
	mlog.Close()
}
