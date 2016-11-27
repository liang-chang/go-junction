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

func make(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	symbolics := conf.Symbolic

	var varDoTarget = makeDoTarget
	var varDoLink = makeDoLink

	errCnt, warnCnt := TraversalSymbolic(&conf, symbolics, varDoTarget, varDoLink)

	tmpl = template.Must(template.New("check_template").Parse(check_template))

	var confMap map[string]interface{} = structs.Map(conf)

	confMap["ErrorCount"] = errCnt
	confMap["WarnCount"] = warnCnt

	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}
}

func makeDoTarget(target string, symbolic *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	if symbolic.Skip {
		return
	}
	ret, _ := util.DirectoryExist(target)
	if ret == false {
		if !conf.Config.CreateTargetFolder {
			symbolic.Target = `Error! "` + target + `"  not exist ! `
			errCnt = 1
			return
		}

		err := os.MkdirAll(target, os.ModePerm)
		if err != nil {
			symbolic.Target = `Error! "` + target + `"  create fail ! ` + err.Error()
			errCnt = 1
		}
	}
	return
}

func makeDoLink(target string, folderIndex int, linkConfig *config.LinkConfig, symb *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	if folderIndex == -1 {
		var warnText = ""
		if !linkConfig.WarnIgnore {
			warnCnt++
			warnText = "Warning !"
		}
		linkConfig.MatchFolder = append(linkConfig.MatchFolder, warnText + " No directory match !")
		return
	}

	if symb.Skip {
		return
	}

	link := linkConfig.MatchFolder[folderIndex]

	var ret bool
	ret, _ = util.Exist(link)
	//link 文件夹不存在，直接创建
	if ret == false {
		var err = os.MkdirAll(link, os.ModePerm)
		if err != nil {
			symb.Target = `Error! "` + link + `"  create fail ! ` + err.Error()
			errCnt = 1
			return
		}
		var success bool
		success, err = symbolic.CreateJunction(link, target, true);
		if !success {
			symb.Target = `Error! "` + link + `"  create junction fail ! ` + err.Error()
			errCnt = 1
			return
		}
		return
	}

	if conf.Config.BackupLinkFolder || linkConfig.Backup {
		var err error
		//TODO 判断之前是否已有备份
		if err = os.Rename(link, link + FOLDER_BACK_SUBFFIX); err != nil {
			linkConfig.MatchFolder[folderIndex] = `Error ! ` + link + `   -->   backup "` + link + `" failed ! ` + err.Error()
			errCnt = 1
			return
		}
	}

	ret, _ = util.IsReparsePoint(link)
	if ret == false {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  not is not junction point !`
		errCnt = 1
		return
	}

	t, err := symbolic.GetJunctionTarget(link)
	if err != nil && t == "" {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  not is not junction point !`
		errCnt = 1
		return
	}

	t = strings.Replace(strings.ToLower(t), `\`, "/", -1)
	if t != strings.ToLower(target) {
		linkConfig.MatchFolder[folderIndex] = `Error ! "` + link + `"  link to : ` + t
		errCnt = 1
		return
	}

	//if

	return
}



