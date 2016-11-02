package main

import (
	"directory"
	"config"
	"os"
	"fmt"
	//"junction"
	//"syscall"
)

func main() {
	fmt.Println(os.Symlink("v:/temp", "v:/aa"))
	config := config.Parse();

	fmt.Println(config.PathAlias)
	fmt.Println(config.Junction)

	ret := directory.GetPatternDir(`d:/|\d+$|/bin`)
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
	fmt.Println(os.RemoveAll("v:/tt"))
	fmt.Println(os.RemoveAll("v:/aa"))
	//fmt.Println(directory.DirectoryExist("v:/xxxxx.txt"))
	//fmt.Println(junction.GetJunctionTarget("V:/tt"));
	//fmt.Println(junction.IsJunction("v:/tt"))
	//fmt.Println(junction.Delete("v:/ttt"))
	//fmt.Println(junction.CreateJunction("v:/tt", "v:/temp",true))
}

