package main

import (
	"github.com/scale-it/go-log"
	"os"
)

var Log = log.New()

func main() {
	Log.AddHandler(os.Stderr, log.Levels.Info, log.StdFormatter{"[root]", log.Lshortfile, false})
	Log.AddHandler(os.Stderr, log.Levels.Info, log.StdFormatter{"[with time]", log.Lshortfile | log.Ldate | log.Ltime, false})
	x := 1.23
	Log.Error("x is %v", x)
	Log.Debug("This won't be printed")

	Log = log.NewStd(os.Stderr, log.Levels.Debug, "", log.Ldate|log.Lmicroseconds, true)
	Log.Debug("some message")
	Log.Info("some message")
	Log.Warning("some message")
	Log.Error("some message")
	Log.Fatal("some message")
}
