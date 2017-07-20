package util

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func (l *Logger) log(i ...interface{}) {
	l.logger.Print(i)
}

var logger Logger

func init() {
	logger = Logger{}
	logger.logger = log.New(os.Stdout, "", log.Ltime)
}
