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
	fmt.Println(os.Symlink("v:/temp","v:/aa"))
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

	//fileInfo, err := os.Stat("v:/tt");

	//fmt.Println(junction.GetJunctionTarget("V:/tt"));
	//fmt.Println(junction.IsJunction("v:/tt"))
	//fmt.Println(junction.Delete("v:/ttt"))
	fmt.Println(junction.Create("h:/aaa", "v:/temp",true))
}

