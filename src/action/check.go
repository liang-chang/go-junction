package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
	//"fmt"
	"util"
	"fmt"
)

func check(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	symbolics := conf.Symbolic

	for sidex, symboT := range conf.Symbolic {

		var target string = symboT.Target

		ret, _ := util.DirectoryExist(target)
		if ret == false {
			symbolics[sidex].Target = "ERROR  ->  " + target +"  ------->  Folder Not exist!"
		}

		var linkConfigs []config.LinkConfig = conf.Symbolic[sidex].LinkConfig

		//symboMap := symbolics[sidex]

		for lindex, _ := range linkConfigs {

			linkConf := linkConfigs[lindex]

			fmt.Println(linkConf)

			var matchFolder []string = linkConfigs[lindex].MatchFolder

			for mi, folder := range matchFolder {
				ret, _ := util.DirectoryExist(folder)
				if ret == false {
					matchFolder[mi] = "NE -> " + folder
				}
			}

		}

	}

	tmpl = template.Must(template.New("check_template").Parse(check_template))

	var confMap map[string]interface{} = structs.Map(conf)

	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}
}



