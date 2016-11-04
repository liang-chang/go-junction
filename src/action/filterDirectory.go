package action

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

func GetPatternDirectory(pathPattern string) []string {

	diskReg := regexp.MustCompile(`(?i)[a-z]:`)

	matchRet := diskReg.FindStringSubmatch(pathPattern)

	if len(matchRet) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	var diskName string = matchRet[0]

	dirReg := regexp.MustCompile(`(\\|/)((?:\|[^|]+\|)|[^/\\]+)`)

	pathPattern = strings.Trim(pathPattern, " ")

	var parsedDirs []string = make([]string, 0, 16)

	//路径部分
	dirParts := dirReg.FindAllStringSubmatch(pathPattern, -1)

	if len(dirParts) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	parsedDirs = append(parsedDirs, diskName)

	dirPartsLen := len(dirParts)

	for index, v := range dirParts {

		fileSplit := v[1]

		namePattern := v[2]

		ingoreExist := index + 1 == dirPartsLen

		//普通路径字串
		//最后一个路径可以不存在
		if !(strings.HasPrefix(namePattern, "|") && strings.HasSuffix(namePattern, "|")) {
			parsedDirs = appendDir(parsedDirs, fileSplit, namePattern, ingoreExist)
			continue
		}

		//包含正则的路径字串
		namePattern = strings.Trim(namePattern, "|")

		parsedDirs = appendPatternDir(parsedDirs, fileSplit, namePattern, ingoreExist)
	}
	return parsedDirs
}

/**
dirs 中添加 appendDir
*/
func appendPatternDir(parsedDirs []string, fileSplit string, namePattern string, ingoreExist bool) []string {

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
func appendDir(parsedDirs []string, fileSplit string, appendDir string, ingoreExist bool) []string {
	match := func(d os.FileInfo) bool {
		return d.IsDir()&&d.Name() == appendDir
	}

	retDirs := traversalDir(parsedDirs, fileSplit, match);
	return retDirs
}

type dirMatch func(os.FileInfo) bool

func traversalDir(parsedDirs []string, fileSplit string, match dirMatch) []string {
	var retDirs []string = make([]string, 0, 16)
	for _, v := range parsedDirs {
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
