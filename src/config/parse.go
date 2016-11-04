package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const (
	actionParamName = "action"

	actionDefault = "list" //list check create recovery

	actoinUsageText = "action can only one of those: check create update"

	//参数默认名称
	configParamName = "config"

	//默认配置文件夹名
	configFileDefaultName = "config.toml"

	//提示语
	configUsageText = "config file name"
)

func Parse() config {
	var actionName = flag.String(actionParamName, actionDefault, actoinUsageText)

	var configFileName = flag.String(configParamName, configFileDefaultName, configUsageText)

	flag.Parse();

	var config config

	if _, err := toml.DecodeFile(*configFileName, &config); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	config.ActionName = *actionName

	return config
}