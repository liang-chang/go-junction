package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	//"fmt"
	"util"
)

func check(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	symbolics := conf.Symbolic

	var doTarget = func(target string, symbolic *config.Symbolic) (errCnt int) {
		ret, _ := util.DirectoryExist(target)
		if ret == false {
			symbolic.Target = "ERROR  ->  " + target + "  ------->  Folder Not exist!"
			errCnt = 1
		} else {
			errCnt = 0
		}
		return errCnt
	}

	var doLink = func(target, link string, folderIndex int, linkConfig *config.LinkConfig) (errCnt int) {
		ret, _ := util.DirectoryExist(link)
		if ret == false {
			linkConfig.MatchFolder[folderIndex] = "NE -> " + link
			errCnt = 1
		}
		errCnt = 0
		return errCnt
	}

	errCnt, warnCnt := TraversalSymbolic(symbolics, doTarget, doLink)

	tmpl = template.Must(template.New("check_template").Parse(check_template))

	var confMap map[string]interface{} = structs.Map(conf)

	confMap["ErrorCount"] = errCnt
	confMap["WarnCount"] = warnCnt



	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}
}



