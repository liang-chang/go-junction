package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	"util"
	"symbolic"
)

func check(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	symbolics := conf.Symbolic

	var varDoTarget = checkDoTarget
	var vardoLink = checkDoLink

	errCnt, warnCnt := TraversalSymbolic(&conf, symbolics, varDoTarget, vardoLink)

	tmpl = template.Must(template.New("check_template").Parse(check_template))

	var confMap map[string]interface{} = structs.Map(conf)

	confMap["ErrorCount"] = errCnt
	confMap["WarnCount"] = warnCnt

	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}
}

func checkDoTarget(target string, symbolic *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	ret, _ := util.DirectoryExist(target)
	if ret == false {
		symbolic.Target = `Error! "` + target + `"  not exist !`
		errCnt = 1
	} else {
		errCnt = 0
	}
	return
}

func checkDoLink(target string, index int, linkConf *config.LinkConfig, symb *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	if index == -1 {
		var warnText = ""
		if !linkConf.WarnIgnore {
			warnCnt++
			warnText = "Warning !"
		}
		linkConf.MatchFolder = append(linkConf.MatchFolder, warnText + " No directory match !")
		return
	}

	link := linkConf.MatchFolder[index]
	var ret bool
	var msg = link
	var oldRealTarget string

	ret, _ = util.Exist(link)
	if ret == false {
		msg = `Error ! "` + link + `" not exist !`
		goto setErrorText
	}

	ret, _ = util.IsReparsePoint(link)
	if ret == false {
		msg = `Error ! "` + link + `"  is not junction point !`
		goto setErrorText
	}

	oldRealTarget, _ = symbolic.GetJunctionTarget(link)

	if oldRealTarget == "" {
		msg = `Error ! "` + link + `"  is not junction point !`
		goto setErrorText
	}

	if linkConf.Isolate {
		if !util.IsSubDirectory(oldRealTarget, target) {
			msg = `Error ! "` + link + `"  link to : ` + oldRealTarget
			goto setErrorText
		}

	} else {
		if !util.IsSamePath(oldRealTarget, target) {
			msg = `Error ! "` + link + `"  link to : ` + oldRealTarget
			goto setErrorText
		}
	}

	linkConf.MatchFolder[index] = msg
	return

	setErrorText:
	linkConf.MatchFolder[index] = msg
	errCnt = 1
	return
}



