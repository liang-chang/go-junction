package action

import (
	"config"
	"fmt"
	"strconv"
	"strings"
	"os"
	"text/template"
)

func list(conf config.Setting) {
	t := template.Must(template.ParseFiles("common.tmpl"))
	t.Execute(os.Stdout, nil)
}

func listBack(conf config.Setting) {
	fmt.Println("Action : " + conf.Action)
	fmt.Println("BackupLinkFolder  : " + strconv.FormatBool(conf.Config.BackupLinkFolder))
	fmt.Println("ClearBackupFolder : " + strconv.FormatBool(conf.Config.ClearBackupFolder))
	fmt.Println("CreateTargetFolder: " + strconv.FormatBool(conf.Config.CreateTargetFolder))

	for _, symbo := range conf.Symbolic {
		fmt.Println(strings.Repeat("-", 15) + "symbolic" + strings.Repeat("-", 15))
		fmt.Println("action : " + symbo.Action)
		fmt.Println("target : " + symbo.Target)

		for j, linkConfig := range symbo.LinkConfig {
			fmt.Println()
			fmt.Println("link   : " + symbo.Link[j])
			atIndex := strings.Index(symbo.Link[j], "@")
			if atIndex < 0 {
				atIndex = 0
			} else {
				atIndex++
			}
			width := strconv.Itoa(8 + atIndex)
			for k, folder := range linkConfig.MatchFolder {
				if k == 0 {
					fmt.Printf("match  :%" + strconv.Itoa(atIndex) + "s %s\n", "", folder)
					continue
				}
				fmt.Printf("%" + width + "s %s\n", "", folder)
			}

		}
	}
}

