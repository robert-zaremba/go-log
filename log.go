/* This implements an alternative logger to the one found in the standard
 * library with support for more logging levels, formatters and outputs.
 * The main goal is to provide easy and flexible way to handle new outputs and formats
 * Author: Robert Zaremba
 *
 * https://github.com/scale-it/go-log
 */
package log

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// Represents how critical the logged
// message is.
type Level uint8

var Levels = struct {
	Trace   Level
	Debug   Level
	Info    Level
	Warning Level
	Error   Level
	Fatal   Level
}{0, 10, 20, 30, 40, 50}

// Verbose names of the levels
var levelStrings = map[Level]string{
	Levels.Trace:   "TRACE",
	Levels.Debug:   "DEBUG",
	Levels.Info:    "INFO",
	Levels.Warning: "WARN",
	Levels.Error:   "ERROR",
	Levels.Fatal:   "FATAL",
}

// Verbose and colored names of the levels
var levelCStrings = map[Level]string{
	Levels.Trace:   levelStrings[Levels.Trace],
	Levels.Debug:   levelStrings[Levels.Debug],
	Levels.Info:    AnsiEscape(MAGENTA, levelStrings[Levels.Info], OFF),
	Levels.Warning: AnsiEscape(YELLOW, levelStrings[Levels.Warning], OFF),
	Levels.Error:   AnsiEscape(RED, levelStrings[Levels.Error], OFF),
	Levels.Fatal:   AnsiEscape(RED, BOLD, levelStrings[Levels.Fatal], OFF),
}

// Returns an log Level which name match given string.
// If there is no such Level, then Levels.Debug is returned
func String2Level(level string) (Level, error) {
	if level == "" {
		return Levels.Debug, errors.New("level is empty")
	}
	for li, ls := range levelStrings {
		if ls == level {
			return li, nil
		}
	}
	return Levels.Debug, errors.New("Wrong log level " + level)
}

type output struct {
	writer io.Writer
	level  Level
}

// The Logger
type Logger struct {
	mtx     sync.Mutex
	outputs []output
}

// Instantiate a new Logger
func New() *Logger {
	return &Logger{sync.Mutex{}, make([]output, 0)}
}

// Adds an ouput, specifying the maximum log Level
// you want to be written to this output. For instance,
// if you pass Warning for level, all logs of type
// Warning, Error, and Fatal would be logged to this output.
func (this *Logger) AddOutput(writer io.Writer, level Level) {
	this.mtx.Lock()
	this.outputs = append(this.outputs, output{writer, level})
	this.mtx.Unlock()
}

// Convenience function
func (this *Logger) Trace(format string, v ...interface{}) {
	// TODO: split the string
	this.Logger(Levels.Trace, format, v...)
}

// Convenience function
func (this *Logger) Debug(format string, v ...interface{}) {
	this.Logger(Levels.Debug, format, v...)
}

// Convenience function
func (this *Logger) Info(format string, v ...interface{}) {
	this.Logger(Levels.Info, format, v...)
}

// Convenience function
func (this *Logger) Warning(format string, v ...interface{}) {
	this.Logger(Levels.Warning, format, v...)
}

// Convenience function
func (this *Logger) Error(format string, v ...interface{}) {
	this.Logger(Levels.Error, format, v...)
}

// Convenience function, will not terminate the program
func (this *Logger) Fatal(format string, v ...interface{}) {
	this.Logger(Levels.Fatal, format, v...)
}

// Loggers a message for the given level. Most callers will likely
// prefer to use one of the provided convenience functions.
func (this *Logger) Logger(level Level, format string, v ...interface{}) {
	message := fmt.Sprintf(format+"\n", v...)
	strTimestamp := getTimestamp()
	strFinal := fmt.Sprintf("%s [%-5s] %s", strTimestamp, levelCStrings[level], message)
	this.log(level, strFinal)
}

func (this *Logger) log(level Level, m string) {
	bytes := []byte(m)
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, output := range this.outputs {
		if output.level <= level {
			output.writer.Write(bytes)
		}
	}
}

// Gets the timestamp string
func getTimestamp() string {
	now := time.Now()
	return fmt.Sprintf("%v-%02d-%02d %02d:%02d:%02d.%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000)
}
