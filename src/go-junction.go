package main

import (
	"directory"
	"config"
	//	"fmt"
	//	"regexp"
	//"log"
	//	"os"
	//	"os/exec"
	//	"path/filepath"
	//	"strings"
	//	"bytes"
	//	"strings"
	//	"bytes"
	//	"path/filepath"
	//	"os"
	//	"log"
	//"io/ioutil"
	//"fmt"

	"os"
	"fmt"
	"junction"
	//"syscall"
)

func main() {

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

	//fmt.Println(junction.GetTarget("V:/TEMP - 目录连接点"));
	fmt.Println(junction.GetTarget("V:/tt"));
	//fmt.Println(junction.GetTarget("d:/Users/ZL/AppData/Roaming/Tencent/QQ/temp"));
	fmt.Println(junction.Exists("v:/tt"))
	//fmt.Println(junction.Delete("v:/ttt"))
	//fmt.Println(junction.Create("v:/aaa", "v:/temp",true))
}

