package log_test

import (
	"fmt"
	"github.com/scale-it/go-log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var formatter = log.SimpleFormatter{}
var fn = "rotedfile.tmp"

func TestRotFileNoRot(t *testing.T) {
	var l *log.Logger = log.New(fmt.Sprint, fmt.Sprintf)
	var err error
	rotfile, err := log.NewRotFile(fn, true, 5, 0)
	defer os.Remove(fn)
	if err != nil {
		t.Fatal(err)
	}
	l.AddHandler(&rotfile, log.Levels.Debug, formatter)
	for i := 1; i < 5; i++ {
		l.Log(log.Levels.Debug, "1")
	}
	time.Sleep(time.Millisecond)
	assert_content(t, "1111", fn)
}

func TestRotFileRotations(t *testing.T) {
	var l *log.Logger = log.New(fmt.Sprint, fmt.Sprintf)
	rotfile, _ := log.NewRotFile(fn, true, 5, 2)
	defer os.Remove(fn)
	defer os.Remove(fn + ".1")
	defer os.Remove(fn + ".2")
	l.AddHandler(&rotfile, log.Levels.Debug, formatter)
	for i := 1; i <= 17; i++ {
		l.Log(log.Levels.Debug, "1")
	}
	time.Sleep(time.Millisecond)
	assert_content(t, "11", fn)
	assert_content(t, "11111", fn+".1")
	assert_content(t, "11111", fn+".2")
}

func LogALot(ch chan bool, rotfile *log.RotFile, s string) {
	var l *log.Logger = log.New(fmt.Sprint, fmt.Sprintf)
	l.AddHandler(rotfile, log.Levels.Debug, formatter)
	for i := 1; i <= 2000; i++ {
		l.Log(log.Levels.Debug, s)
		// to assure a gorutine switch
		if i%500 == 0 {
			time.Sleep(time.Microsecond)
		}
	}
	ch <- true
}

func TestRotFileGoroutines(t *testing.T) {
	var i uint
	maxBakcups := uint(6)
	rotfile, _ := log.NewRotFile(fn, true, 5000, maxBakcups)
	ch := make(chan bool)
	for i = 0; i < 10; i++ {
		go LogALot(ch, &rotfile, strconv.Itoa(int(i))+" ")
	}
	for i = 0; i < 10; i++ { // wait for log workers to finish
		<-ch
	}

	dir := filepath.Dir(fn)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("Error while reading files from directory %s; %v", dir, err)
	}
	re, _ := regexp.Compile(fn + `\.\d+`)
	i = 0
	found_base := false
	for _, fi := range files {
		name := fi.Name()
		if !fi.IsDir() && re.MatchString(name) {
			i++
			os.Remove(name)
		}
		if name == fn {
			found_base = true
			os.Remove(name)
		}
	}
	if !found_base {
		t.Errorf("'%s' log base file doesn't found", fn)
	}
	if i != maxBakcups {
		t.Errorf("There should be %d backups instad of %d", maxBakcups, i)
	}
}

func assert_content(t *testing.T, content string, filename string) {
	var err error
	var n int
	var f *os.File
	clen := len(content)
	if f, err = os.Open(filename); err != nil {
		t.Fatal(err)
	}
	var log_content = make([]byte, clen+1)
	if n, err = f.Read(log_content); err != nil {
		t.Fatal(err)
	}
	if n != clen {
		t.Fatalf("writed %d bytes instead of %d", n, clen)
	}
	if string(log_content[:clen]) != content {
		t.Fatalf("Wrong log content. Should be '%s', but is %s", content, log_content)
	}
}
