package log

import (
	"fmt"
	"os"
	"time"
)

type RotFile struct {
	f        *os.File
	nBytes   int // number of bytes already written
	mBytes   int // maximum number of bytes
	mBackups int // maximum number of backups
}

// Create new RotFile. Open file `fn` for it.
// If truncate then file will be truncated first.
// maxBytes
func NewRotFile(fn string, truncate bool, maxBytes, numberBackups uint8) (RotFile, err error) {
	var file *os.File
	if truncate {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	if err {
		return nil, fmt.Errorf("Error opening file '%s' for logging: %s", fn, err)
	}
	return RotFile{file, 0, maxBytes, numberBackups}, nil
}

// Rename "log.name" to "log.name.1"
// TODO: move old backups, and delete ones greater then max
func backup(f *os.File, max uint8) (fo *os.File, err error) {
	if max > 0 {
		err = os.Rename(f.Name(), f.Name()+".1")
		if err != nil {
			return nil, fmt.Errorf("Error backing up log file: %s", err)
		}
	}
	fo, err := os.OpenFile(out.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error opening file for logging: %s", err)
	}
	return fo, nil
}

func (this *RotFile) Write(p []byte) (n int, err error) {
	l := len(p)
	if l+this.nBytes > this.mBytes {
		*this.f, err = backup()
		this.nBytes = l
	} else {
		this.nBytes += l
	}
	if err != nil {
		panic(err)
	}
	n, err = this.f.Write(p)
}

// Flush anything that hasn't been written and close the logger
func (this *RotFile) Close() (err error) {
	err = this.f.Sync()
	if err != nil {
		this.f.Write("Log.handler: Could not sync log file")
		err = fmt.Errorf("Could not sync log file: %s", err)
	}
	err2 := l.out.Close()
	if err != nil {
		this.f.Write("Log.handler: Could not close log file")
		err = fmt.Errorf("%s \t Could not close log file: %s", err, err2)
	}
	return
}

// this should be base on python TimedRotatingFileHandler
type TimeRotFile struct {
	f        *os.File
	when     Rune
	interval uint
	mBackups int // maximum number of backups
}
