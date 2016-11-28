package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	"util"
	"symbolic"
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

	//link 文件夹不存在，直接创建
	if ret, _ := util.Exist(link); !ret {
		goto createJunction
	}

	//如果需要备份
	if linkConfig.Backup || conf.Config.BackupLinkFolder {
		//如果之前已经有备份
		if ret, _ := util.Exist(link + FOLDER_BACK_SUBFFIX); ret {
			//删除现有文件夹的内容
			if err := util.RemoveContents(link); err != nil {
				linkConfig.MatchFolder[folderIndex] = `Error !  can not remove "` + link + `" content ! ` + err.Error()
				errCnt = 1
				return
			}
			goto createJunction
		} else {
			//如果之前没有备份

			//如果需要删除备份文件夹的内容
			if linkConfig.Clear || conf.Config.ClearBackupFolder {
				if err := util.RemoveContents(link); err != nil {
					linkConfig.MatchFolder[folderIndex] = `Error !  can not remove "` + link + `" content ! ` + err.Error()
					errCnt = 1
					return
				}
			}

			//重命名备命
			if err := os.Rename(link, link + FOLDER_BACK_SUBFFIX); err != nil {
				linkConfig.MatchFolder[folderIndex] = `Error ! backup failed ! can not rename "` + link + `" ! ` + err.Error()
				errCnt = 1
				return
			}
		}

	}

	createJunction:
	if err := os.MkdirAll(link, os.ModePerm); err != nil {
		symb.Target = `Error! directory "` + link + `"  create fail ! ` + err.Error()
		errCnt = 1
		return
	}
	ret, err := symbolic.CreateJunction(link, target, true);
	if err != nil || !ret {
		symb.Target = `Error! "` + link + `"  create junction fail ! ` + err.Error()
		errCnt = 1
		return
	}

	return
}



