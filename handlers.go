package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	//"time"
)

type RotFile struct {
	f        *os.File
	nBytes   int   // number of bytes already written
	mBytes   int   // maximum number of bytes
	mBackups uint8 // maximum number of backups
}

// Create new RotFile. Open file `fn` for it.
// If truncate then file will be truncated first.
// maxBytes
func NewRotFile(fn string, truncate bool, maxBytes int, numberBackups uint8) (RotFile, error) {
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
func (this *RotFile) Close() (err error) {
	err = this.f.Sync()
	if err != nil {
		this.f.Write([]byte("Log.handler: Could not sync log file"))
		err = fmt.Errorf("Could not sync log file: %s", err)
	}
	err2 := this.f.Close()
	if err != nil {
		this.f.Write([]byte("Log.handler: Could not close log file"))
		err = fmt.Errorf("%s \t Could not close log file: %s", err, err2)
	}
	return
}

/* Roll files:
   - remove "target.<x>" for x >=max
   - rename "target.<x>" to "target.<x+1>"
   - rename "target" to "target.1"   */
func rotFiles(target string, max int) (err error) {
	base := filepath.Base(target)
	dir := filepath.Dir(target)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Error reading files from directory %s; %v", base, err)
	}
	re, _ := regexp.Compile(`tt\.` + `\d+`)
	var i int
	var to_rollI = make(map[int]string)
	for _, fi := range files {
		name := fi.Name()
		if !fi.IsDir() && re.MatchString(name) {
			roll_num := name[strings.LastIndex(name, ".")+1:]
			i, _ = strconv.Atoi(roll_num)
			if i >= max { // remove unnecesary backups
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
		println(k)
	}
	sort.Ints(to_roll)

	// roll files in descending order
	for i := len(to_roll) - 1; i >= 0; i-- {
		num := to_roll[i]
		name := to_rollI[num]
		println("renaming ", name)
		if err = os.Rename(name, fmt.Sprintf("%s.%d", base, num+1)); err != nil {
			return fmt.Errorf("Error backing up log file: %s", err)
		}
	}

	if err = os.Rename(base, base+".1"); err != nil {
		return fmt.Errorf("Error backing up log file: %s", err)
	}
	return nil
}

func (this *RotFile) Write(p []byte) (n int, err error) {
	l := len(p)
	if l+this.nBytes > this.mBytes {
		if err = this.Close(); err != nil {
			println(err)
			this.f.Write([]byte(err.Error()))
		} else {
			err = rotFiles(this.f.Name(), this.mBytes)
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
