package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

func GetMatchDirectory(pathPattern string, linkConfg LinkConfig) []string {

	diskReg := regexp.MustCompile(`(?i)[a-z]:`)//`(?i)[a-z]:(\\|/)`

	matchRet := diskReg.FindStringSubmatch(pathPattern)

	if len(matchRet) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	var diskName string = matchRet[0]

	dirReg := regexp.MustCompile(`(?:\\|/)((?:\|[^|]+\|)|[^/\\]+)`)//`(\\|/)((?:\|[^|]+\|)|[^/\\]+)`

	pathPattern = strings.Trim(pathPattern, " ")

	var matchedDirs []string = make([]string, 0, 16)

	//路径部分
	dirParts := dirReg.FindAllStringSubmatch(pathPattern, -1)

	if len(dirParts) < 1 {
		fmt.Println(pathPattern + " 路径无效")
		return make([]string, 0, 0)
	}

	matchedDirs = append(matchedDirs, diskName + FILE_SPLIT)

	dirPartsLen := len(dirParts)

	for index, v := range dirParts {

		namePattern := v[1]

		isLast := index + 1 == dirPartsLen

		//普通路径字串
		//最后一个路径可以不存在
		if !(strings.HasPrefix(namePattern, "|") && strings.HasSuffix(namePattern, "|")) {
			matchedDirs = appendDir(matchedDirs, namePattern, linkConfg, isLast)
			continue
		}

		//包含正则的路径字串
		namePattern = strings.Trim(namePattern, "|")
		//最后一路径是正则的必须要存在
		matchedDirs = appendPatternDir(matchedDirs, namePattern, linkConfg, isLast)
	}
	return matchedDirs
}

/**
dirs 中添加 appendDir
*/
func appendPatternDir(parsedDirs []string, namePattern string, linkConfg, isLast bool) []string {

	nameReg := regexp.MustCompile(namePattern)

	match := func(d os.FileInfo) bool {
		return d.IsDir()&&nameReg.MatchString(d.Name())
	}

	retDirs := traversalDir(parsedDirs, namePattern, match, linkConfg);

	return retDirs
}

/**
dirs 中添加 appendDir
*/
func appendDir(parsedDirs []string, appendDir string, linkConfg LinkConfig, isLast bool) []string {
	if linkConfg.ForeceCreate {
		var ret []string = make([]string, 0, 16)
		for _, v := range parsedDirs {
			//fix：出现 e://temp 这种路径情况
			v = trimFileSplit(v)


		}
	}

	match := func(d os.FileInfo) bool {
		return d.IsDir()&&d.Name() == appendDir
	}

	retDirs := traversalDir(parsedDirs, appendDir, match, linkConfg);
	return retDirs
}

type dirMatch func(os.FileInfo) bool

func traversalDir(parsedDirs []string, appendDir string, match dirMatch, linkConfg LinkConfig) []string {
	var retDirs []string = make([]string, 0, 16)
	for _, v := range parsedDirs {

		childDirs, _ := ioutil.ReadDir(v)

		//fix：出现 e://temp 这种路径情况
		v = trimFileSplit(v)

		for _, d := range childDirs {
			if match(d) {
				retDirs = append(retDirs, v + FILE_SPLIT + d.Name())
			}
		}
	}
	return retDirs
}

func trimFileSplit(str string) string {
	if strings.HasSuffix(str, FILE_SPLIT) {
		return strings.TrimRight(str, FILE_SPLIT)
	}
	return str
}