package main

import (
	"github.com/scale-it/go-log"
	"os"
)

var Log = log.New()

func main() {
	//Log.AddOutput(os.Stdout, clog.LevelWarning)
	Log.AddOutput(os.Stderr, log.Levels.Info)
	x := 1.23
	Log.Info("x is %v", x)
	Log.Debug("This won't be printed")
}
