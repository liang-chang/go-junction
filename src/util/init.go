package util

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Ltime)
}

func Log(i ...interface{}) {
	logger.Println(i...)
}

func Logf(fmt string, i ...interface{}) {
	logger.Printf(fmt, i...)
}
