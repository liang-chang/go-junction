package main

import (
	//"util"
	"config"
	"os"
	"fmt"
	//"junction"
	//"syscall"
	//"action"
	//"util"
	//"unsafe"
	//	"symbolic"
)

func main() {

	configs := config.Init();

	fmt.Println(configs.PathAlias)
	//fmt.Println(configs.Symbolic)
	symbolics := configs.Symbolic

	fmt.Println("-----------------")
	fmt.Println(configs.Symbolic)

	for _, v := range symbolics {
		fmt.Println(v)
		fmt.Println(v.LinkConfig)
	}

	fmt.Println(symbolics)

	linkConfig := config.LinkConfig{};

	ret := config.GetMatchDirectory(`E:/|workspace|/.metadataT`, linkConfig)
	//ret := action.GetPatternDirectory(`E:/|\d+$|/binn`)
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
	//fmt.Println(symbolic.GetJunctionTarget("V:/tt"));
	//fmt.Println(junction.IsJunction("v:/tt"))
	//fmt.Println(symbolic.DeleteJunction("v:/tt"))
	//fmt.Println(symbolic.CreateJunction("v:/tt", "v:/TEMP",true))

	//var mountPoint symbolic.MountPointReparseBuffer
	//var symbolicLink symbolic.SymbolicLinkReparseBuffer
	//
	//fmt.Println(unsafe.Sizeof(mountPoint))
	//fmt.Println(unsafe.Sizeof(symbolicLink))
}

