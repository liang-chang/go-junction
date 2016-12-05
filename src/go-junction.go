package main

import (
	"config"
	"action"
	"log"
	"os"
	"fmt"
	"symbolic"
)

var logger *log.Logger;

func init() {
	logger = log.New(os.Stdout, "", log.Ltime)
}

func main() {

	logger.Printf("Read config ……")
	confSetting := config.Read();

	logger.Printf("MatchDirectory ……")
	config.MatchDirectory(&confSetting)

	logger.Printf("Call action:%s ……", confSetting.Action)
	action.Call(confSetting.Action, confSetting);


	fmt.Println(symbolic.CreateJunction("v:/tt", "v:/temp", true))
}

