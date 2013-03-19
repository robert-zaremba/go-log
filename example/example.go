package main

import (
	"github.com/scale-it/go-log"
	log2 "log"
	"os"
)

var Log = log.New()

func main() {
	//Log.AddOutput(os.Stdout, clog.LevelWarning)
	Log.AddOutput(os.Stderr, log.Levels.Info, log.StdFormatter{"root ", log2.Lshortfile, false})
	Log.AddOutput(os.Stderr, log.Levels.Info, log.StdFormatter{"with time", log2.Lshortfile | log2.Ldate | log2.Ltime, false})
	x := 1.23
	Log.Info("x is %v", x)
	Log.Debug("This won't be printed")

	Log = log.NewStd(os.Stderr, log.Levels.Debug, "", log2.Ldate|log2.Lmicroseconds, true)
	Log.Debug("some message")
	Log.Info("some message")
	Log.Warning("some message")
	Log.Error("some message")
	Log.Fatal("some message")
}
