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
	logger.Print(i)
}
