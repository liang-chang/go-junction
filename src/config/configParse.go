package config

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"util"

	"github.com/BurntSushi/toml"
)

const (
	ACTION_PARAM_NAME = "action"

	ACTION_DEFAULT = "list" //list check make recovery

	ACTOIN_USAGE_TEXT = "action can only one of those: list check make recovery"

	//参数默认名称
	CONFIG_PARAM_NAME = "config"

	//默认配置文件夹名
	CONFIG_FILE_DEFAULT_NAME = "config.toml"

	//提示语
	CONFIG_USAGE_TEXT = "config file name"

	//auto target 的 auto 命令
	AUTO_TARGET_AUTO = `<auto>`

	//auto target 的 next 命令
	AUTO_TARGET_NEXT = `<next>`
)

/*
   读取配置文件，并进行配置的解读；进行 path alias 的替换
*/
func Read() Setting {
	//读取解析配置文件
	setting := readConfig()

	//解析 TargetFolderPattern
	//设置所有 targetFolder
	setTargetFoldersByPattern(&setting)

	setBuildInPathAlias(&setting)

	parseTargetLinkCmd(&setting)

	return setting
}

func parseTargetLinkCmd(setting *Setting) {
	targetFolderSlice := setting.Config.TargetFolders

	//解析 link 和 target 文件夹中的替换符
	//解析 link 字符中的千个命令
	//从可用的target 目录下剔除补手动映射的目录
	for i, symbCopy := range setting.Symbolic {

		symbolic := &setting.Symbolic[i]

		(*symbolic).Target = resolvePathAlias((*symbolic).Target, setting.PathAlias)

		trimTraget := strings.TrimSpace((*symbolic).Target)

		var err error
		if trimTraget != AUTO_TARGET_AUTO && trimTraget != AUTO_TARGET_NEXT {

			if _, err = os.Stat(trimTraget); err == nil {
				//格式化文件路径
				if trimTraget, err = filepath.Abs(trimTraget); err != nil {
					util.Logf("target=%v 是个无效文件夹", trimTraget)
					os.Exit(1)
				}

				found := util.SliceIndex(len(targetFolderSlice), func(i int) bool { return targetFolderSlice[i] == trimTraget })

				if found >= 0 {
					targetFolderSlice = append(targetFolderSlice[:found], targetFolderSlice[found+1:]...)
				}

			}
		}

		for j, linkText := range symbCopy.Link {
			(*symbolic).Link[j] = resolvePathAlias(linkText, setting.PathAlias)
			(*symbolic).LinkConfig = append((*symbolic).LinkConfig, readLinkText((*symbolic).Link[j], setting.PathAlias))
		}
	}

	setting.Config.TargetFolders = targetFolderSlice
}

func setTargetFoldersByPattern(setting *Setting) {

	var tartgetMap = make([]string, 0, 50)
	config := setting.Config
	for _, value := range config.TargetFolderPattern {
		matches, _ := filepath.Glob(value)
		if len(matches) > 0 {
			for _, m := range matches {
				tartgetMap = append(tartgetMap, m)
			}
		}
	}
	setting.Config.TargetFolders = tartgetMap
}

func setBuildInPathAlias(conf *Setting) {
	usr, _ := user.Current()

	conf.PathAlias["UserHome"] = strings.Replace(usr.HomeDir, `\`, `/`, -1)

	conf.PathAlias["Temp"] = strings.Replace(os.TempDir(), `\`, `/`, -1)

	for key, value := range conf.PathAlias {
		conf.PathAlias[key] = resolvePathAlias(value, conf.PathAlias)
	}

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

	ret.FolderPattern = split[0]
	return ret
}

func setLinkCmd(cmd string, linkConfg *LinkConfig) {
	cmd = strings.ToLower(cmd)
	if strings.Contains(cmd, "b") {
		linkConfg.Backup = true
		if strings.Contains(cmd, "c") {
			linkConfg.Clear = true
		}
	}

	if strings.Contains(cmd, "i") {
		linkConfg.Isolate = true
	}

	if strings.Contains(cmd, "l") {
		linkConfg.LastDirAppender = true
	}
	if strings.Contains(cmd, "f") {
		linkConfg.ForeceCreate = true
	}

	if strings.Contains(cmd, "w") {
		linkConfg.WarnIgnore = true
	}

}

func readConfig() Setting {
	var actionName = flag.String(ACTION_PARAM_NAME, ACTION_DEFAULT, ACTOIN_USAGE_TEXT)

	var configFileName = flag.String(CONFIG_PARAM_NAME, CONFIG_FILE_DEFAULT_NAME, CONFIG_USAGE_TEXT)

	flag.Parse()

	//设置工作目录
	workDirPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(os.Args[0]))
	syscall.SetCurrentDirectory(workDirPtr)

	configExist, _ := util.FileExist(*configFileName)

	if !configExist {
		log.Fatal(*configFileName + " file not exist!")
		os.Exit(1)
	}

	var conf = Setting{}

	if _, err := toml.DecodeFile(*configFileName, &conf); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	valid, _ := util.Contain(Actions, *actionName)

	if !valid {
		log.Fatal("unknown action " + *actionName)
		os.Exit(1)
	}

	conf.ConfigFile = *configFileName

	conf.Action = *actionName

	return conf
}
