package main

import (
	"config"
	//"action"
	//"strconv"
	//"fmt"
	//"symbolic"
	"action"
	"log"
	"os"
	"fmt"
	"util"
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

	//ret := action.GetPatternDirectory(`E:/|\d+$|/binn`)

	//fmt.Println(directory.DirectoryExist("v:/tt"))
	//fmt.Println(directory.DirectoryExist("v:/xxxxx"))

	//os.RemoveAll 可以删了作 符号链接和 junction point
	//fmt.Println(os.RemoveAll("v:/tt"))
	//fmt.Println(os.RemoveAll("v:/aa"))

	//fmt.Println(util.IsReparsePoint("v:/xxxxx.txt"))
	//fmt.Println(util.IsReparsePoint("v:/temp"))
	//fmt.Println(util.IsReparsePoint("v:/t1"))
	//fmt.Println(util.IsReparsePoint("v:/t2"))
	fmt.Println(util.Exist("v:/useless/aaa"))
	fmt.Println(util.IsReparsePoint("v:/t1"))
	fmt.Println(util.IsReparsePoint("v:/t2"))
	//fmt.Println(symbolic.GetJunctionTarget("V:/tt"));
	//fmt.Println(symbolic.GetJunctionTarget("V:/t2"));
	//fmt.Println(symbolic.IsJunction("v:/tt"))
	//fmt.Println(symbolic.DeleteJunction("v:/tt"))
	//fmt.Println(symbolic.CreateJunction("v:/tt", "V:/TEMP",true))

	//var mountPoint symbolic.MountPointReparseBuffer
	//var symbolicLink symbolic.SymbolicLinkReparseBuffer

	// /fmt.Println(unsafe.Sizeof(mountPoint))
	//fmt.Println(unsafe.Sizeof(symbolicLink))
}

