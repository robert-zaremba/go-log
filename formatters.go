package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Formatter interface {
	Format(Level, string) []byte
}

// Standard Formatter
type StdFormatter struct {
	Prefix  string // prefix to write at beginning of each line
	Flag    int    // format flags - based flags from std log package
	Colored bool   // use colored level names
}

func (this StdFormatter) Format(level Level, msg string) []byte {
	var slevel string
	var ok bool
	var out []string

	// adding time info
	if this.Flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		now := time.Now()
		if this.Flag&Ldate != 0 {
			out = append(out, fmt.Sprintf("%v-%02d-%02d", now.Year(), now.Month(), now.Day()))
		}
		if this.Flag&(Lmicroseconds) != 0 {
			out = append(out, fmt.Sprintf("%02d:%02d:%02d.%06d", now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000))
		} else if this.Flag&(Ltime) != 0 {
			out = append(out, fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second()))
		}
	}

	// adding level info
	if this.Colored {
		slevel, ok = levelCStrings[level]
	} else {
		slevel, ok = levelStrings[level]
	}
	if !ok {
		slevel = strconv.Itoa(int(level))
	}
	out = append(out, slevel)

	// adding caller info. It's quiet exepnsive
	if this.Flag&(Lshortfile|Llongfile) != 0 {
		if _, file, line, ok := runtime.Caller(4); ok { // calldepth, 4 functions back
			if this.Flag&Lshortfile != 0 {
				file = file[strings.LastIndex(file, "/")+1:]
			}
			out = append(out, fmt.Sprintf("%s:%d", file, line))
		} else {
			out = append(out, "???")
		}
	}

	out = append(out, this.Prefix)
	out = append(out, msg)
	return []byte(strings.Join(out, " "))
}

type TimeFormatter struct {
	Prefix string // prefix to write at beginning of each line
}

func (this TimeFormatter) Format(level Level, msg string) []byte {
	var out []string

	now := time.Now().UTC()
	out = append(out, fmt.Sprintf("%v-%02d-%02d %02d:%02d:%02d.%06d UTC", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000))
	out = append(out, this.Prefix)
	out = append(out, msg)
	return []byte(strings.Join(out, " "))
}
