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
)

func main() {
	fmt.Println(os.Args)

	config.Init();

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

}

