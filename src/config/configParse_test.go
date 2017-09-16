package config

import (
	"flag"
	"os"
	"testing"
	"util"
	"path/filepath"
	"io/ioutil"
)

func TestReadConfig(t *testing.T) {
	//--config=config.toml  --action=make

	var configContent = `
		[config]
		# link文件夹重命名备份
		backupLinkFolder = false

		# 清空备份的文件夹
		clearBackupFolder = true

		#当target文件不存在时，创建
		createTargetFolder = true

		#当没有匹配到link文件夹时是否警告，默认为false
		warnNoMatchLinkFolder = false

		#当target文件自动分配
		targetFolders=[
		'v:/useless/Z/Z[0-9]',
		]

		[pathAlias]
		#build in path variable
		# UserHome
		# Temp
		useless='V:/useless/'
		chromeCache='V:/chrome_cache'

		[[symbolic]]
		target = '{useless}/Z'
		link = [
			'fil@v:/|log.|/tt',
		]

		[[symbolic]]
		target = '<auto>'
		link = [
			'bcilf@v:/|log.|/tt',
		]
	`

	var configFile= filepath.Dir(os.Args[0])+"/config.toml"
	os.Remove(configFile)

	ioutil.WriteFile(configFile,[]byte(configContent),0777)

	os.Args[1]="--config=config.toml"
	os.Args[2]="--action=check"

	util.Log(os.Args)

	confSetting :=Read()
	//flag.Var()
	util.Log(confSetting)
}

func TestFlag(t *testing.T) {
	util.Log(os.Args)

	os.Args = append(os.Args, ` --config=ttttt.toml  --action=make`)

	util.Log(os.Args)

	readConfig()

	util.Log(flag.Args())
}
