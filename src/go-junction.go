package main

import (
	"config"
	//"action"
	//"strconv"
	//"fmt"
	//"symbolic"
	//"symbolic"
	//"fmt"
	//"symbolic"
	"action"
	"symbolic"
	"fmt"
)

func main() {

	confSetting := config.Read();

	config.MatchDirectory(&confSetting)

	action.Call(confSetting.Action, confSetting);

	//ret := action.GetPatternDirectory(`E:/|\d+$|/binn`)

	//fmt.Println(directory.DirectoryExist("v:/tt"))
	//fmt.Println(directory.DirectoryExist("v:/xxxxx"))

	//os.RemoveAll 可以删了作 符号链接和 junction point
	//fmt.Println(os.RemoveAll("v:/tt"))
	//fmt.Println(os.RemoveAll("v:/aa"))

	//fmt.Println(directory.DirectoryExist("v:/xxxxx.txt"))
	fmt.Println(symbolic.GetJunctionTarget("V:/tt"));
	fmt.Println(symbolic.GetJunctionTarget("V:/t2"));
	//fmt.Println(symbolic.IsJunction("v:/tt"))
	//fmt.Println(symbolic.DeleteJunction("v:/tt"))
	//fmt.Println(symbolic.CreateJunction("v:/tt", "V:/TEMP",true))

	//var mountPoint symbolic.MountPointReparseBuffer
	//var symbolicLink symbolic.SymbolicLinkReparseBuffer

	// /fmt.Println(unsafe.Sizeof(mountPoint))
	//fmt.Println(unsafe.Sizeof(symbolicLink))
}

