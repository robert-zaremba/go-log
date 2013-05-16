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
	Log.Info("x is", x)
	Log.Infof("x is %f and is positive: %v", x, x > 0)
	Log.Debug("This won't be printed")

	Log = log.NewStd(os.Stderr, log.Levels.Debug, log.Ldate|log.Lmicroseconds|log.Lshortfile, true)
	Log.Log(15, "raw message\n")
	Log.Debug("some message")
	Log.Info("some message")
	Log.Warning("some message")
	Log.Error("some message")
	Log.Critical("some message")
}
