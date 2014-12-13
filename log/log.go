package log

import (
	"log"
	"os"
)

var (
	ERROR      *log.Logger
	WARN       *log.Logger
	DEBUG      *log.Logger
	INFO       *log.Logger
	initalised bool = false
)

func setup() {
	if initalised {
		return
	}
	ERROR = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	WARN = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	INFO = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	DEBUG = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	initalised = true
}

func Error(msg string, args ...interface{}) {
	setup()
	ERROR.Printf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	setup()
	WARN.Printf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	setup()
	INFO.Printf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	setup()
	DEBUG.Printf(msg, args...)
}
