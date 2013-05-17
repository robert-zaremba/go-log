package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// DailyFile is thread safe and can be used by multiple Loggers thanks to
// go-log global mutex map. But the filename can't be used elsewhere for write.
type RotFile struct {
	f        *os.File
	nBytes   int  // number of bytes already written
	mBytes   int  // maximum number of bytes
	mBackups uint // maximum number of backups
}

/* Create new RotFile. Open file `fn` for it. Other function arguments:
   - truncate: if true, then file will be truncated first.
   - maxBytes: number of bytes which must be written to rot file
   - numberBackup: number of backups to keep. */
func NewRotFile(fn string, truncate bool, maxBytes int, numberBackups uint) (RotFile, error) {
	var err error
	var file *os.File
	if truncate {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	if err != nil {
		return RotFile{}, fmt.Errorf("Error opening file '%s' for logging: %s", fn, err)
	}
	return RotFile{file, 0, maxBytes, numberBackups}, nil
}

// Flush anything that hasn't been written and close the logger
func close(f *os.File) (err error) {
	err = f.Sync()
	if err != nil {
		f.Write([]byte("Log.handler: Could not sync log file"))
		err = fmt.Errorf("Could not sync log file: %s", err)
	}
	err2 := f.Close()
	if err != nil {
		f.Write([]byte("Log.handler: Could not close log file"))
		err = fmt.Errorf("%s \t Could not close log file: %s", err, err2)
	}
	return
}

/* Roll files:
   - remove "target.<x>" for x >=max
   - rename "target.<x>" to "target.<x+1>"
   - rename "target" to "target.1"   */
func rotFiles(target string, max uint) (err error) {
	imax := int(max)
	base := filepath.Base(target)
	base_len := len(base)
	dir := filepath.Dir(target)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Error while reading files from directory %s; %v", dir, err)
	}
	re, _ := regexp.Compile(base + `\.\d+`)
	var i int
	var to_rollI = make(map[int]string)
	for _, fi := range files {
		name := fi.Name()
		if !fi.IsDir() && re.MatchString(name) {
			i, _ = strconv.Atoi(name[base_len+1:])
			if i >= imax { // remove unnecesary backups
				os.Remove(filepath.Join(dir, name))
			} else {
				to_rollI[i] = name
			}
		}
	}
	// sort rolling names
	var to_roll []int
	for k := range to_rollI {
		to_roll = append(to_roll, k)
	}
	sort.Ints(to_roll)

	// roll files in descending order
	for i := len(to_roll) - 1; i >= 0; i-- {
		num := to_roll[i]
		name := to_rollI[num]
		if err = os.Rename(name, fmt.Sprintf("%s.%d", base, num+1)); err != nil {
			return fmt.Errorf("Error backing up log file: %s", err)
		}
	}

	// roll base file
	if max > 0 {
		if err = os.Rename(base, base+".1"); err != nil {
			return fmt.Errorf("Error backing up log file: %s", err)
		}
	}
	return nil
}

// When using from Logger, it is protected against simultaneous write in Logger.Log
// Logger Log method (which calls this Write) acquires Mutex for this writer.
func (this *RotFile) Write(p []byte) (n int, err error) {
	l := len(p)
	if l+this.nBytes > this.mBytes {
		if err = close(this.f); err != nil {
			n, _ = this.f.Write([]byte(err.Error()))
			return
		} else {
			err = rotFiles(this.f.Name(), this.mBackups)
			f, err2 := os.OpenFile(this.f.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err2 != nil {
				panic("Can't open file for logging. " + err2.Error())
			}
			if err != nil {
				println(err)
				f.Write([]byte(err.Error()))
			}
			this.nBytes = l
			this.f = f
		}
	} else {
		this.nBytes += l
	}
	return this.f.Write(p)
}

// this should be base on python TimedRotatingFileHandler
type TimeRotFile struct {
	f        *os.File
	when     rune
	interval uint
	mBackups int // maximum number of backups
}
