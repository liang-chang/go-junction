package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
	"util"
	"syscall"
	"strings"
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

func Init() Setting {
	conf := readConfig()

	for _, symb := range conf.Symbolic {
		symb.LinkConfig = make([]LinkConfig, len(symb.Link) * 3)

		for _, linkText := range symb.Link {
			readLinkText(linkText)
		}

	}
	return Setting{}
}

func readLinkText(linkText string) LinkConfig {
	split := strings.Split(linkText, "@")
	if len(split) > 2 {
		log.Fatal("unknown text" + linkText)
		os.Exit(1)
	}

	ret := LinkConfig{}


	//解析形如： bclf@d:/|\d+$|/bin
	//前面的 bclf
	if len(split) > 1 {
		cmd := split[0]
		setLinkCmd(cmd, ret)
		split = split[1:]
	}

	return LinkConfig{}
}

func setLinkCmd(cmd string, linkConfg LinkConfig) {
	if strings.Contains(cmd, "b") {
		linkConfg.Backup = true
		if strings.Contains(cmd, "c") {
			linkConfg.Clear = true
		}
	}
	if strings.Contains(cmd, "l") {
		linkConfg.LastDirAppender = true
	}
	if strings.Contains(cmd, "f") {
		linkConfg.ForeceCreate = true
	}

}

func readConfig() Setting {
	var actionName = flag.String(actionParamName, actionDefault, actoinUsageText)

	var configFileName = flag.String(configParamName, configFileDefaultName, configUsageText)

	flag.Parse();

	//设置工作目录
	workDirPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(os.Args[0]));
	syscall.SetCurrentDirectory(workDirPtr)

	configExist, _ := util.FileExist(*configFileName)

	if !configExist {
		log.Fatal(*configFileName + " file not exist!")
		os.Exit(1)
	}

	var conf Setting

	if _, err := toml.DecodeFile(*configFileName, &conf); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	valid, _ := util.Contain(Actions, *actionName)

	if !valid {
		log.Fatal("unknown action " + *actionName)
		os.Exit(1)
	}

	conf.Action = *actionName

	return conf
}