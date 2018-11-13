package main

import (
	"fmt"
	"io"
	"log"
)

const (
	LogLevelNull    = 0
	LogLevelTrace   = 1
	LogLevelDebug   = 2
	LogLevelInfo    = 3
	LogLevelWarning = 4
	LogLevelError   = 5
	LogLevelFatal   = 6
)

var defaultLogLevel uint8 = LogLevelDebug

type Logger struct {
	*log.Logger
}

func (this *Logger) SetFlags(flag int) {
	log.SetFlags(flag)
}

func (this *Logger) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func (this *Logger) SetLevel(level uint8) {
	defaultLogLevel = level
}

func (this *Logger) Trace(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelTrace {
		return
	}
	log.Output(2, string("[TRACE] ")+fmt.Sprintf(format, v...))
}

func (this *Logger) Debug(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelDebug {
		return
	}
	log.Output(2, string("[DEBUG] ")+fmt.Sprintf(format, v...))
}

func (this *Logger) Info(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelInfo {
		return
	}
	log.Output(2, string("[INFO] ")+fmt.Sprintf(format, v...))
}

func (this *Logger) Warning(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelWarning {
		return
	}
	log.Output(2, string("[WARNING] ")+fmt.Sprintf(format, v...))
}

func (this *Logger) Error(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelError {
		return
	}
	log.Output(2, string("[ERROR] ")+fmt.Sprintf(format, v...))
}

func (this *Logger) Fatal(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelFatal {
		return
	}
	log.Output(2, string("[FATAL] ")+fmt.Sprintf(format, v...))
}