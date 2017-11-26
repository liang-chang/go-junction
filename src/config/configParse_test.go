package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"util"
)

func TestMain(m *testing.M) {

	workDirPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(os.Args[0]))
	syscall.SetCurrentDirectory(workDirPtr)

	//单元测试前创文件夹
	var cache = []string{"chrome", "firefox", "opera", "safari"}
	for _, name := range cache {
		os.Mkdir(name+"_cache", os.ModePerm)
	}

	var start = int('A')
	var end = int('Z')

	for i := start; i <= end; i++ {
		os.MkdirAll(`useless/`+string(rune(i)), os.ModePerm)
	}

	os.Exit(m.Run())
}

func TestReadConfig(tt *testing.T) {
	//--config=config.toml  --action=make

	var preSetting = Setting{
		Action:     "check",
		ConfigFile: "config_test.toml"}

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
		'useless/?',
		'*_cache',
		]

		[pathAlias]
		#build in path variable
		# UserHome
		# Temp
		useless='useless'
		chromeCache='chrome_cache'
		tempCache='{Temp}/cache'

		[[symbolic]]
		target = '{useless}/ZE'
		link = [
			'fil@v:/|log.|/tt',
		]

		[[symbolic]]
		target = '{useless}/A'
		link = [
			'bcilf@{UserHome}/tt',
		]

		[[symbolic]]
		target = 'safari_cache'
		link = [
			'bcilf@{UserHome}/tt',
		]
	`

	var configFile = filepath.Dir(os.Args[0]) + "/" + preSetting.ConfigFile

	os.Remove(configFile)

	ioutil.WriteFile(configFile, []byte(configContent), 0777)

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

	prev := reflect.ValueOf(*preSetting.Config)

	pt := reflect.TypeOf(*parsedConfig)
	pv := reflect.ValueOf(*parsedConfig)
	for i := 0; i < pv.NumField(); i++ {
		if pv.Field(i).CanInterface() { //判断是否为可导出字段

			var isEqual bool
			switch pv.Field(i).Kind() {
			case reflect.Bool:
				isEqual = pv.Field(i).Interface() == prev.Field(i).Interface()
			case reflect.String:
				isEqual = pv.Field(i).Interface() == prev.Field(i).Interface()
			case reflect.Array:
				isEqual = reflect.DeepEqual(pv.Field(i).Interface(), prev.Field(i).Interface())
			case reflect.Slice:
				isEqual = reflect.DeepEqual(pv.Field(i).Interface(), prev.Field(i).Interface())
			case reflect.Map:
				isEqual = reflect.DeepEqual(pv.Field(i).Interface(), prev.Field(i).Interface())
			default:
				isEqual = false
			}
			fmt.Printf("%v : equal = %v \n", pt.Field(i).Name, isEqual)
			if !isEqual {
				tt.Errorf("%v != %v \n", pv.Field(i).Interface(), prev.Field(i).Interface())
			}
			//util.Logf("%s %s = %v \n",
			//	pt.Field(i).Name,
			//	pt.Field(i).Type,
			//	pv.Field(i).Interface())
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

//有一个无效的symbolic，以及一个文夹件被批定
//os.exit 1
func TestReadConfig2(tt *testing.T) {
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

		PathAlias: make(map[string]string),
		Symbolic: []*Symbolic{
			&Symbolic{
				Target: `{useless}/ZE`,
				Link:   []string{`fil@v:/|log.|/tt`},
			},
			&Symbolic{
				Target: `'{useless}/A'`,
				Link:   []string{`bcilf@{UserHome}/tt`},
			},
			&Symbolic{
				Target: `v:/safari_cache`,
				Link:   []string{`bcilf@{UserHome}/tt`},
			},
		},
	}

	//解析 TargetFolderPattern
	//设置所有 targetFolder
	setTargetFoldersByPattern(&preSetting)

	setBuildInPathAlias(&preSetting)

	parseTargetLinkCmd(&preSetting)

	util.Log("Action=", preSetting.Action)

	util.Log("ConfigFile=", preSetting.ConfigFile)

	parsedConfig := preSetting.Config

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

	util.Log("PathAlias=", preSetting.PathAlias)

	for _, s := range preSetting.Symbolic {
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

func TestFlag(t *testing.T) {
	var preSetting = Setting{
		Action:     "check",
		ConfigFile: "config_test.toml"}

	targetFolderPattern := [1]string{"v:/useless/?"}

	var preGlobalConfig = GlobalConfig{
		BackupLinkFolder:      false,
		ClearBackupFolder:     true,
		CreateTargetFolder:    false,
		TargetFolderPattern:   targetFolderPattern[:],
		WarnNoMatchLinkFolder: false,
	}

	preSetting.Config = &preGlobalConfig

	fmt.Printf("1=%v \n", preSetting.Config)

	fmt.Printf("1=%v \n", preGlobalConfig)

	preGlobalConfig.TargetFolderPattern = nil

	fmt.Printf("2=%v \n", preSetting.Config)

	fmt.Printf("2=%v \n", preGlobalConfig)
}

func TestSliceRemove(t *testing.T) {
	arr := [5]int{0, 1, 2, 3, 4}
	slice := make([]int, 0, 50)
	slice = append(slice, arr[:]...)

	index := 0
	slice = append(slice[:index], slice[index+1:]...)
	fmt.Print(slice)

	index = 3
	slice = append(slice[:index], slice[index+1:]...)
	fmt.Print(slice)
}

func TestFilepathGlob(t *testing.T) {
	matches, err := filepath.Glob("V:/useless/?")

	util.Log(matches)
	util.Log(err)

	util.Logf("%s %s = %v \n", "a", "b", "c")

	fmt.Print(fmt.Sprintf("%s %s = %v \n", "a", "b", "c"))

	s, _ := filepath.Abs("V:/useless//Z")
	fmt.Printf(s)
}
