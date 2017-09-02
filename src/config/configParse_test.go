package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"util"
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
		firefoxCache='V:/firefox_cache'
		appDataRoaming='{UserHome}/AppData/Roaming'
		appDataLocal='{UserHome}/AppData/Local'
		QQDataHome='{UserHome}/AppData/Roaming/Tencent'

		[[symbolic]]
		target = '@'
		link = [
		'filw@v:/|log.|/tt',
		]
	`
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//tmpfile, _ := ioutil.TempFile(dir, "config_test")

	//flag.Var()

	util.Log(configContent)
	util.Log(dir)

	//flag.Var()
	util.Log("a")
}

func TestFlag(t *testing.T) {
	util.Log(os.Args)

	os.Args = append(os.Args, ` --config=ttttt.toml  --action=make`)

	util.Log(os.Args)

	readConfig()

	util.Log(flag.Args())
}
