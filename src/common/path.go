package common

import (
	"fmt"
	"regexp"
	"strings"
)

func GetPatternDir(path string) {
	dirReg := regexp.MustCompile(".*?(\\|/|$)")
	var parsePath []string = make([]string, 5)
	for {
		loc := dirReg.FindStringSubmatchIndex(path)
		curDir := path[loc[0]:loc[1]]
		curDir = strings.Trim(curDir, "")
		//普通路径字串
		if !(strings.HasPrefix(curDir, "|") && strings.HasSuffix("|")) {
			append(parsePath, curDir)
		}

		//包含正则的路径字串
		strings.Trim(curDir, "|")

		break
	}

	fmt.Println(parsePath)
}
