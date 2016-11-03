package action

import (
	"util"
	"config"
	"os"
	"fmt"
	//"junction"
	//"syscall"
	"symbolic"
)

func Traversal() {
	config := config.Parse();

	fmt.Println(config.PathAlias)
	fmt.Println(config.Junction)

	ret := util.GetPatternDir(`d:/|\d+$|/bin`)
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
	//fmt.Println(junction.Delete("v:/ttt"))
	fmt.Println(symbolic.CreateJunction("v:/tt", "v:/TEMP",true))
}

