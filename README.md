Log
===

Lightweight, alternative Go logger with elastic formatter and multiple output support.


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
	Log.Debug("some message")
	Log.Info("some message")
	Log.Warning("some message")
	Log.Error("some message")
	Log.Fatal("some message")
}
```

Check *example* directory for more.
