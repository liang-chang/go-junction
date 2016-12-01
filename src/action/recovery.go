package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	"util"
)

func recovery(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	symbolics := conf.Symbolic

	var varDoTarget = recoveryDoTarget

	var varDoLink = recoveryDoLink

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

func recoveryDoTarget(target string, symbolic *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
	//恢复 target 不用做什么处理
	return
}

func recoveryDoLink(target string, folderIndex int, linkConfig *config.LinkConfig, symb *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int) {
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

	//如果之前已经有备份
	if ret, _ := util.Exist(link + FOLDER_BACK_SUBFFIX); ret {
		if err := os.RemoveAll(link); err != nil {
			linkConfig.MatchFolder[folderIndex] = `Error ! can not remove symbo link "` + link + `" ! ` + err.Error()
			errCnt = 1
			return
		}

		//备份还原
		if err := os.Rename(link + FOLDER_BACK_SUBFFIX, link); err != nil {
			linkConfig.MatchFolder[folderIndex] = `Error ! backup failed ! can not rename "` + link + `" ! ` + err.Error()
			errCnt = 1
			return
		}
		return
	}


	//link 文件夹不存在，直接创建
	if ret, _ := util.Exist(link); !ret {
		goto recovery
		return
	} else {
		//如果是符号链接或者是 junction 直接删除
		if isReparsePoint, _ := util.IsReparsePoint(link); isReparsePoint {
			if err := os.RemoveAll(link); err != nil {
				linkConfig.MatchFolder[folderIndex] = `Error ! can not remove symbo link "` + link + `" ! ` + err.Error()
				errCnt = 1
				return
			}
			goto recovery
		}
		return
	}

	recovery:

	if err := os.MkdirAll(link, os.ModePerm); err != nil {
		linkConfig.MatchFolder[folderIndex] = ` Error! directory "` + link + `"  create fail ! ` + err.Error()
		errCnt = 1
		return
	}
	return
}



