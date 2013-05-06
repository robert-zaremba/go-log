Log
===

Lightweight, alternative Go logger with elastic formatter and multiple output support.

Main features:

* thread safe
* support colored output
* logging levels: Debug, Info, Warning, Error, Fatal


Example
-------

```go
package main

import (
	"github.com/scale-it/go-log"
	"os"
)

func main() {
	Log := log.NewStd(os.Stderr, log.Levels.Debug, log.Ldate|log.Lmicroseconds, true)
	Log.Log(10, "non parametrized message") // level, message
	param := 1.23
	Log.Debug("some message, %d", param)
	Log.Info("some message, %d", param)
	Log.Warning("some message, %d", param)
	Log.Error("some message, %d", param)
	Log.Fatal("some message, %d", param)
}
```

Check *example* directory for more.
