package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	exitFunc func(int)
}

func New() *Logger {
	return &Logger{
		Logger: log.New(
			os.Stdout,
			"[homebox-export] ",
			log.LstdFlags,
		),
		exitFunc: os.Exit,
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		l.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Printf("[FATAL] "+format, v...)
	l.exitFunc(1)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
}
