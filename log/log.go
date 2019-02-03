package log

import (
	"fmt"
	"log"
	"os"
)

var std = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

func SetLogger(logger *log.Logger) {
	std = logger
}

type Logger struct {
	Id string
}

func (l *Logger) Printf(format string, v ...interface{}) {
	format = l.Id + ": " + format
	std.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	format = l.Id + ": " + format
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}