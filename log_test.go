package log_test

import (
	"fmt"
	"github.com/scale-it/go-log"
	"os"
)

func Example() {
	var Log *log.Logger = log.New(fmt.Sprintln, fmt.Sprintf)
	formatter := log.StdFormatter{"[root]", log.Lshortfile, false}
	Log.AddHandler(os.Stdout, log.Levels.Warning, formatter)
	pi := 3.14
	Log.Infof("Pi is %v", pi)
	Log.Warningf("Pi is %v", pi)
	Log.Error("Pi is", pi)
	// Output:
	// WARN  log_test.go:15 [root] Pi is 3.14
	// ERROR log_test.go:16 [root] Pi is 3.14
}
