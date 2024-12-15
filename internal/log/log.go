package log

import (
	"strings"

	"github.com/labstack/gommon/log"
)

type Logger struct {
	*log.Logger
}

var logger *Logger

func InitLogger(level string) {
	l := log.New("-")

	lev := log.DEBUG
	switch strings.ToUpper(level) {
	case "DEBUG":
		lev = log.DEBUG
	case "INFO":
		lev = log.INFO
	case "ERROR":
		lev = log.ERROR
	case "WARN":
		lev = log.WARN
	}

	l.EnableColor()
	l.SetHeader(`{"time":"${time_rfc3339_nano}","level":"${level}","file":"${short_file}:${line}"}`)
	l.SetLevel(lev)

	logger = &Logger{Logger: l}
}

func GetLogger() *Logger {
	return logger
}

func (l *Logger) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}
