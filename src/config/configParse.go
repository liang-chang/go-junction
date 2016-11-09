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
	"os/user"
	"regexp"
	//"fmt"
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

/*
     读取配置文件，并进行配置的解读；进行 path alias 的替换
 */
func Read() Setting {
	conf := readConfig()
	setBuildInPathAlias(conf)
	for i, symb := range conf.Symbolic {
		for _, linkText := range symb.Link {
			conf.Symbolic[i].LinkConfig = append(conf.Symbolic[i].LinkConfig, readLinkText(linkText, conf.PathAlias))
		}
		conf.Symbolic[i].Target = resolvePathAlias(conf.Symbolic[i].Target, conf.PathAlias)
	}
	return conf
}

func setBuildInPathAlias(conf Setting) {
	usr, _ := user.Current()

	conf.PathAlias["UserHome"] = usr.HomeDir

	conf.PathAlias["Temp"] = os.TempDir()
}

func resolvePathAlias(folderPattern string, pathAlias map[string]string) string {
	enablelias := regexp.MustCompile(`{[^}]+}`).FindAllStringSubmatch(folderPattern, -1)
	resovelPath := folderPattern
	for _, v := range enablelias {
		aliasName := strings.TrimRight(strings.TrimLeft(v[0], "{"), "}")
		aliasVal, ok := pathAlias[aliasName]
		if !ok {
			log.Fatal("unknown path alias : " + aliasName + " in " + folderPattern)
		}
		resovelPath = strings.Replace(folderPattern, v[0], aliasVal, -1)
	}
	return resovelPath
}

func readLinkText(linkText string, pathAlias map[string]string) LinkConfig {
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
		setLinkCmd(cmd, &ret)
		split = split[1:]
	}

	ret.FolderPattern = resolvePathAlias(split[0], pathAlias);
	return ret
}

func setLinkCmd(cmd string, linkConfg *LinkConfig) {
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