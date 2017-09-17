package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"util"
)

func TestReadConfig(tt *testing.T) {
	//--config=config.toml  --action=make

	var configContent = `
		[config]
		# link文件夹重命名备份
		backupLinkFolder = false

		# 清空备份的文件夹
		clearBackupFolder = true

		#当target文件不存在时，创建
		createTargetFolder = false

		#当没有匹配到link文件夹时是否警告，默认为false
		warnNoMatchLinkFolder = false

		#当target文件自动分配
		TargetFolderPattern=[
		'v:/useless/?',
		]

		[pathAlias]
		#build in path variable
		# UserHome
		# Temp
		useless='V:/useless/'
		chromeCache='V:/chrome_cache'
		tempCache='{Temp}/cache'

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

	var configFile = filepath.Dir(os.Args[0]) + "/config.toml"
	os.Remove(configFile)

	ioutil.WriteFile(configFile, []byte(configContent), 0777)

	os.Args[1] = "--config=config.toml"
	os.Args[2] = "--action=check"

	util.Log(os.Args)

	confSetting := Read()

	util.Log("Action=", confSetting.Action)

	util.Log("ConfigFile=", confSetting.ConfigFile)

	config := confSetting.Config

	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() { //判断是否为可导出字段
			util.Logf("%s %s = %v \n",
				t.Field(i).Name,
				t.Field(i).Type,
				v.Field(i).Interface())
		}
	}

	util.Log("PathAlias=", confSetting.PathAlias)

	for _, s := range confSetting.Symbolic {
		t := reflect.TypeOf(s)
		v := reflect.ValueOf(s)
		var str = "Symbolic "
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanInterface() { //判断是否为可导出字段
				str += fmt.Sprintf("%s = %v ,",
					t.Field(i).Name,
					v.Field(i).Interface())
			}
		}
		util.Log(str)
	}

}

func TestFlag(t *testing.T) {
	util.Log(os.Args)

	os.Args = append(os.Args, ` --config=ttttt.toml  --action=make`)

	util.Log(os.Args)

	readConfig()

	util.Log(flag.Args())
}

func TestFilepathGlob(t *testing.T) {
	matches, err := filepath.Glob("V:/useless/?")

	util.Log(matches)
	util.Log(err)

	util.Logf("%s %s = %v \n", "a", "b", "c")

	fmt.Print(fmt.Sprintf("%s %s = %v \n", "a", "b", "c"))
}
