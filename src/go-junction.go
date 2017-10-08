package main

import (
	"action"
	"config"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Ltime)
}

func main() {

	logger.Printf("Parse config ……")
	confSetting := config.Read()

	logger.Printf("Match Link And Target Directory ……")
	config.MatchDirectory(&confSetting)

	logger.Printf("Call action:%s ……", confSetting.Action)
	action.Call(confSetting.Action, confSetting)
}
