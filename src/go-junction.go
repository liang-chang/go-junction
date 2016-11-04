package main

import (
	//"util"
	"config"
	"os"
	"fmt"
	//"junction"
	//"syscall"
	"symbolic"
	"action"
	//"unsafe"
)

func main() {
	config := config.Parse();

	fmt.Println(config.PathAlias)
	fmt.Println(config.Junction)

	ret := action.GetPatternDirectory(`d:/|\d+$|/binn`)
	for _, v := range ret {
		_, err := os.Stat(v)
		if err != nil {
			fmt.Println(v + " err")
		}
		if os.IsNotExist(err) {
			fmt.Println(v + " invalid")
		} else {
			fmt.Println(v)
		}
	}

	//fmt.Println(directory.DirectoryExist("v:/tt"))
	//fmt.Println(directory.DirectoryExist("v:/xxxxx"))

	//os.RemoveAll 可以删了作 符号链接和 junction point
	//fmt.Println(os.RemoveAll("v:/tt"))
	//fmt.Println(os.RemoveAll("v:/aa"))

	//fmt.Println(directory.DirectoryExist("v:/xxxxx.txt"))
	//fmt.Println(junction.GetJunctionTarget("V:/tt"));
	//fmt.Println(junction.IsJunction("v:/tt"))
	fmt.Println(symbolic.DeleteJunction("v:/tt"))
	//fmt.Println(symbolic.CreateJunction("v:/tt", "v:/TEMP",true))

	//var mountPoint symbolic.MountPointReparseBuffer
	//var symbolicLink symbolic.SymbolicLinkReparseBuffer
	//
	//fmt.Println(unsafe.Sizeof(mountPoint))
	//fmt.Println(unsafe.Sizeof(symbolicLink))
}

