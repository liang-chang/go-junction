package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	"fmt"
	"util"
)

func check(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	var confMap map[string]interface{} = structs.Map(conf)

	var symbolics []map[string]interface{} = confMap["Symbolic"].([]map[string]interface{});

	for sidex, symbo := range symbolics {

		var target string = symbo["Target"].(string)

		var linkConfigs []config.LinkConfig = conf.Symbolic[sidex].LinkConfig

		symboMap := symbolics[sidex]

		for lindex, linkConfCopy := range linkConfigs {

			linkConf := linkConfigs[lindex]

			fmt.Println(linkConf)

			var matchFolder []string = linkConfCopy["MatchFolder"]

			for folder := range matchFolder {
				ret, _ := util.DirectoryExist(folder)
				if ret == false {

				}
			}

		}

	}

	tmpl = template.Must(template.New("check_template").Parse(check_template))

	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}
}



