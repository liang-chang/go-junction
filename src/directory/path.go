package directory

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

func GetPatternDir(path string)[]string  {

	diskReg := regexp.MustCompile(`(?i)[a-z]:(\\|/)`)

	matchRet := diskReg.FindStringSubmatch(path)

	if len(matchRet) < 1 {
		fmt.Println(path + " 路径无效")
		return make([]string, 0, 0)
	}

	var diskName string = matchRet[0]

	dirReg := regexp.MustCompile(`(\\|/)((?:\|[^|]+\|)|[^/\\]+)`)

	path = strings.Trim(path, " ")

	var parsedDirs []string = make([]string, 0, 16)

	dirNames := dirReg.FindAllStringSubmatch(path, -1)

	if len(dirNames) < 1 {
		fmt.Println(path + " 路径无效")
		return make([]string, 0, 0)
	}

	parsedDirs = append(parsedDirs, diskName)

	for _, v := range dirNames {

		//		path := v[0]

		fileSplit := v[1]

		namePattern := v[2]

		//普通路径字串
		if !(strings.HasPrefix(namePattern, "|") && strings.HasSuffix(namePattern, "|")) {

			parsedDirs=appendDir(parsedDirs, fileSplit, namePattern)
			continue
		}

		//包含正则的路径字串
		namePattern = strings.Trim(namePattern, "|")

		parsedDirs = appendPatternDir(parsedDirs, fileSplit, namePattern)
	}
	return parsedDirs
}

/**
dirs 中添加 appendDir
*/
func appendPatternDir(parsedDirs []string, fileSplit string, namePattern string) []string {

	nameReg := regexp.MustCompile(namePattern)

	match := func(d os.FileInfo) bool {
		return d.IsDir()&&nameReg.MatchString(d.Name())
	}

	retDirs := traversalDir(parsedDirs, fileSplit, match);

	return retDirs
}

/**
dirs 中添加 appendDir
*/
func appendDir(parsedDirs []string, fileSplit string, appendDir string) []string {
	match := func(d os.FileInfo) bool {
		return d.IsDir()&&d.Name() == appendDir
	}

	retDirs := traversalDir(parsedDirs, fileSplit, match);
	return retDirs
}

type dirMatch func(os.FileInfo) bool

func traversalDir(dirs []string, fileSplit string, match dirMatch) []string {
	var retDirs []string = make([]string, 0, 16)
	for _, v := range dirs {
		childDirs, err := ioutil.ReadDir(v)

		if err != nil {
			fmt.Println(v + " 读取子文件夹出错！")
			continue
		}

		for _, d := range childDirs {

			if match(d) {
				retDirs = append(retDirs, v + fileSplit + d.Name())
			}
		}
	}
	return retDirs
}
