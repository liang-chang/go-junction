package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"util"
)

//有一个无效的symbolic，以及一个文夹件被批定
//os.exit 1
func TestMatchDirectory(tt *testing.T) {
	//--config=config.toml  --action=make

	var preSetting = Setting{
		Action:     "check",
		ConfigFile: "config_test.toml",

		Config: &GlobalConfig{
			BackupLinkFolder:      false,
			ClearBackupFolder:     true,
			CreateTargetFolder:    false,
			TargetFolderPattern:   []string{"v:/useless/?", "v:/*_cache"},
			WarnNoMatchLinkFolder: false,
		},
	}

	setTargetFoldersByPattern(&preSetting)

	var configFile = filepath.Dir(os.Args[0]) + "/" + preSetting.ConfigFile

	os.Remove(configFile)

	configContent, _ := json.Marshal(preSetting)

	ioutil.WriteFile(configFile, configContent, 0777)

	var actionName = "check"

	preSetting.Action = actionName

	os.Args[1] = "--config=" + preSetting.ConfigFile
	os.Args[2] = "--action=" + preSetting.Action

	util.Log(os.Args)

	parsedConfSetting := Read()

	util.Log("Action=", parsedConfSetting.Action)

	if parsedConfSetting.Action != preSetting.Action &&
		parsedConfSetting.ConfigFile != preSetting.ConfigFile {

		tt.Fatal("action != %v ,ConfigFile !=%v ", preSetting.Action, parsedConfSetting.ConfigFile)
	}

	util.Log("ConfigFile=", parsedConfSetting.ConfigFile)

	parsedConfig := parsedConfSetting.Config

	pt := reflect.TypeOf(*parsedConfig)
	pv := reflect.ValueOf(*parsedConfig)
	for i := 0; i < pv.NumField(); i++ {
		if pv.Field(i).CanInterface() { //判断是否为可导出字段
			util.Logf("%s %s = %v \n",
				pt.Field(i).Name,
				pt.Field(i).Type,
				pv.Field(i).Interface())
		}
	}

	util.Log("PathAlias=", parsedConfSetting.PathAlias)

	for _, s := range parsedConfSetting.Symbolic {
		t := reflect.TypeOf(*s)
		v := reflect.ValueOf(*s)
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
