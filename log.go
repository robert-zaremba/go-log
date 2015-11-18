/* This implements an alternative logger to the one found in the standard
 * library with support for more logging levels, formatters and handlers.
 * The main goal is to provide easy and flexible way to handle new handlers and formats
 * Author: Robert Zaremba
 *
 * https://github.com/scale-it/go-log
 */
package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Map of mutexes which will protect Writeres for simultaneous write
var writerMutexMap = make(map[io.Writer]*sync.Mutex)

// Represents how critical the logged
// message is.
type Level uint8

var Levels = struct {
	Trace    Level
	Debug    Level
	Info     Level
	Warning  Level
	Error    Level
	Critical Level
}{0, 10, 20, 30, 40, 50}

// Verbose names of the levels
var levelStrings = map[Level]string{
	Levels.Trace:    "TRACE",
	Levels.Debug:    "DEBUG",
	Levels.Info:     "INFO ",
	Levels.Warning:  "WARN ",
	Levels.Error:    "ERROR",
	Levels.Critical: "CRITIC",
}

// Verbose and colored names of the levels
var levelCStrings = map[Level]string{
	Levels.Trace:    levelStrings[Levels.Trace],
	Levels.Debug:    levelStrings[Levels.Debug],
	Levels.Info:     AnsiEscape(MAGENTA, levelStrings[Levels.Info], OFF),
	Levels.Warning:  AnsiEscape(YELLOW, levelStrings[Levels.Warning], OFF),
	Levels.Error:    AnsiEscape(RED, levelStrings[Levels.Error], OFF),
	Levels.Critical: AnsiEscape(RED, BOLD, levelStrings[Levels.Critical], OFF),
}

// Returns an log Level which name match given string.
// If there is no such Level, then Levels.Debug is returned
func String2Level(level string) (Level, error) {
	if level == "" {
		return Levels.Debug, errors.New("level is empty")
	}
	for li, ls := range levelStrings {
		if strings.HasPrefix(ls, level) {
			return li, nil
		}
	}
	return Levels.Debug, errors.New("Wrong log level " + level)
}

type handler struct {
	writer io.Writer
	level  Level
	fmt    Formatter
}

type SprinterT func(a ...interface{}) string
type SprinterFT func(format string, a ...interface{}) string

type Logger struct {
	// Mutex to protect simultaneous appends to handlers
	mtx      sync.Mutex
	handlers []handler
	print    SprinterT
	printf   SprinterFT
}

// Instantiate a new Logger
// It requires two arguments - functions which dumps a variables to a string.
// They should be compatybile with fmt.Sprint* functions. Good options are functions from
// `fmt` or `spew` package.
func New(p SprinterT, pf SprinterFT) *Logger {
	return &Logger{sync.Mutex{}, make([]handler, 0), p, pf}
}

// Convenience function to create logger with StdFormatter
// Uses for
func NewStd(w io.Writer, level Level, flag int, colored bool) *Logger {
	l := New(fmt.Sprintln, fmt.Sprintf)
	l.AddHandler(w, level, StdFormatter{"", flag, colored})
	return l
}

/* LOGGER
 * ------
 */

// Adds a handler, specifying the maximum log Level you want to be written to this output.
// For instance, if you pass Warning for level, all logs of type
// Warning, Error, and Critical would be logged to this handler.
// This method is thread safe. You can use it in multiple goroutines.
// You can also use the same writer in multiple Loggers.
func (this *Logger) AddHandler(writer io.Writer, level Level, fm Formatter) {
	this.mtx.Lock()
	if _, ok := writerMutexMap[writer]; !ok {
		writerMutexMap[writer] = &sync.Mutex{}
	}
	this.handlers = append(this.handlers, handler{writer, level, fm})
	this.mtx.Unlock()
}

// Logs a message for the given level. Most callers will likely
// prefer to use one of the provided convenience functions (Debug, Info...).
// The message is evaluated only when needed - if none handler accept that level,
// message won't be dumped and formatted
func (this *Logger) Log(level Level, v ...interface{}) {
	var notDumped = true
	var msg string
	var out []byte
	for _, h := range this.handlers {
		if h.level <= level {
			if notDumped {
				notDumped = true
				msg = this.print(v...)
			}
			out = h.fmt.Format(level, msg)
			mtx, _ := writerMutexMap[h.writer]
			mtx.Lock()
			h.writer.Write(out)
			mtx.Unlock()
		}
	}
}

// Logs a formatted message message for the given level.
// Wrapper around Log method
// The message is evaluated only when needed - if none handler accept that level,
// message won't be dumped and formatted
func (this *Logger) Logf(level Level, format string, v ...interface{}) {
	var notDumped = true
	var msg string
	var out []byte
	for _, h := range this.handlers {
		if h.level <= level {
			if notDumped {
				notDumped = true
				msg = this.printf(format+"\n", v...)
			}
			out = h.fmt.Format(level, msg)
			mtx, _ := writerMutexMap[h.writer]
			mtx.Lock()
			h.writer.Write(out)
			mtx.Unlock()
		}
	}
}

// Convenience function
func (this *Logger) Trace(v ...interface{}) {
	this.Log(Levels.Trace, v...)
}

// Convenience function
func (this *Logger) Tracef(format string, v ...interface{}) {
	this.Logf(Levels.Trace, format, v...)
}

// Convenience function
func (this *Logger) Debug(v ...interface{}) {
	this.Log(Levels.Debug, v...)
}

// Convenience function
func (this *Logger) Debugf(format string, v ...interface{}) {
	this.Logf(Levels.Debug, format, v...)
}

// Convenience function
func (this *Logger) Info(v ...interface{}) {
	this.Log(Levels.Info, v...)
}

// Convenience function
func (this *Logger) Infof(format string, v ...interface{}) {
	this.Logf(Levels.Info, format, v...)
}

// Convenience function
func (this *Logger) Warning(v ...interface{}) {
	this.Log(Levels.Warning, v...)
}

// Convenience function
func (this *Logger) Warningf(format string, v ...interface{}) {
	this.Logf(Levels.Warning, format, v...)
}

// Convenience function, short version of Warning
func (this *Logger) Warn(v ...interface{}) {
	this.Log(Levels.Warning, v...)
}

// Convenience function, short version of Warningf
func (this *Logger) Warnf(format string, v ...interface{}) {
	this.Logf(Levels.Warning, format, v...)
}

// Convenience function
func (this *Logger) Error(v ...interface{}) {
	this.Log(Levels.Error, v...)
}

// Convenience function
func (this *Logger) Errorf(format string, v ...interface{}) {
	this.Logf(Levels.Error, format, v...)
}

// Convenience function, will not terminate the program
func (this *Logger) Critical(v ...interface{}) {
	this.Log(Levels.Critical, v...)
}

// Convenience function, will not terminate the program
func (this *Logger) Criticalf(format string, v ...interface{}) {
	this.Logf(Levels.Critical, format, v...)
}

// Convenience function, will terminate the program
func (this *Logger) Fatal(v ...interface{}) {
	this.Log(Levels.Critical, v...)
	os.Exit(1)
}

// Convenience function, will terminate the program
func (this *Logger) Fatalf(format string, v ...interface{}) {
	this.Logf(Levels.Critical, format, v...)
	os.Exit(1)
}

// Convinience function to support io.Writer interface
func (this *Logger) Write(p []byte) (n int, err error) {
	n = len(p)
	this.Log(0, string(p))
	return
}
