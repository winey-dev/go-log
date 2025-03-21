/*
Package log provides a simple logging package for Go.

example:

import (
	"github.com/winey-dev/go-log"
)

func main() {
	mlog, err := log.NewLogger("my-app")
	if err != nil {
		panic(err)
	}
	defer mlog.Close()

	mlog.Info("Hello, World!\n")
	mlog.Debug("Hello, World!\n")
	mlog.Warn("Hello, World!\n")
	mlog.Error("Hello, World!\n")
}

*/

package log
