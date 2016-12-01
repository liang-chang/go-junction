package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	"util"
	"symbolic"
	"strings"
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

func checkDoLink(target string, folderIndex int, linkConfig *config.LinkConfig, symb *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	if folderIndex == -1 {
		var warnText = ""
		if !linkConfig.WarnIgnore {
			warnCnt++
			warnText = "Warning !"
		}
		linkConfig.MatchFolder = append(linkConfig.MatchFolder, warnText + " No directory match !")
		return
	}

	link := linkConfig.MatchFolder[folderIndex]
	var ret bool
	ret, _ = util.Exist(link)
	if ret == false {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `" not exist !`
		errCnt = 1
		return
	}

	ret, _ = util.IsReparsePoint(link)
	if ret == false {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  is not junction point !`
		errCnt = 1
		return
	}

	t, err := symbolic.GetJunctionTarget(link)
	if err != nil && t == "" {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  is not junction point !`
		errCnt = 1
		return
	}

	t = strings.Replace(strings.ToLower(t), `\`, "/", -1)
	if t != strings.ToLower(target) {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  link to : ` + t
		errCnt = 1
		return
	}

	return
}



