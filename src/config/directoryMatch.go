package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

func GetPatternDirectory(pathPattern string, linkConfg LinkConfig) []string {

	diskReg := regexp.MustCompile(`(?i)[a-z]:(\\|/)`)

	matchRet := diskReg.FindStringSubmatch(pathPattern)

	if len(matchRet) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	var diskName string = matchRet[0]

	dirReg := regexp.MustCompile(`(\\|/)((?:\|[^|]+\|)|[^/\\]+)`)

	pathPattern = strings.Trim(pathPattern, " ")

	var matchedDirs []string = make([]string, 0, 16)

	//路径部分
	dirParts := dirReg.FindAllStringSubmatch(pathPattern, -1)

	if len(dirParts) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	matchedDirs = append(matchedDirs, diskName)

	dirPartsLen := len(dirParts)

	for _, v := range dirParts {

		fileSplit := v[1]

		namePattern := v[2]

		isLast := index + 1 == dirPartsLen

		//普通路径字串
		//最后一个路径可以不存在
		if !(strings.HasPrefix(namePattern, "|") && strings.HasSuffix(namePattern, "|")) {
			matchedDirs = appendDir(matchedDirs, fileSplit, namePattern, isLast)
			continue
		}


		//包含正则的路径字串
		namePattern = strings.Trim(namePattern, "|")
		//最后一路径是正则的必须要存在
		matchedDirs = appendPatternDir(matchedDirs, fileSplit, namePattern, isLast)
	}
	return matchedDirs
}

/**
dirs 中添加 appendDir
*/
func appendPatternDir(parsedDirs []string, fileSplit string, namePattern string, isLast bool) []string {

	nameReg := regexp.MustCompile(namePattern)

	match := func(d os.FileInfo) bool {
		return d.IsDir()&&nameReg.MatchString(d.Name())
	}

	retDirs := traversalDir(parsedDirs, fileSplit, namePattern, match, isLast);

	return retDirs
}

/**
dirs 中添加 appendDir
*/
func appendDir(parsedDirs []string, fileSplit string, appendDir string, isLast bool) []string {
	match := func(d os.FileInfo) bool {

		return d.IsDir()&&d.Name() == appendDir
	}

	retDirs := traversalDir(parsedDirs, fileSplit, appendDir, match, isLast);
	return retDirs
}

type dirMatch func(os.FileInfo) bool

func traversalDir(parsedDirs []string, fileSplit string, appendDir string, match dirMatch, isLast bool) []string {
	var retDirs []string = make([]string, 0, 16)
	for _, v := range parsedDirs {

		if isLast {
			//fix：出现 e://temp 这种路径情况
			if strings.HasSuffix(v, fileSplit) {
				v = strings.TrimRight(v, fileSplit)
			}
			retDirs = append(retDirs, v + fileSplit + appendDir)
			continue
		}

		childDirs, _ := ioutil.ReadDir(v)

		//fix：出现 e://temp 这种路径情况
		if strings.HasSuffix(v, fileSplit) {
			v = strings.TrimRight(v, fileSplit)
		}

		for _, d := range childDirs {

			if match(d) {

				retDirs = append(retDirs, v + fileSplit + d.Name())
			}
		}
	}
	return retDirs
}
