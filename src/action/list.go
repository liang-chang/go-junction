package action

import (
	"config"
	"os"
	"text/template"
	"github.com/fatih/structs"
)

func list(conf config.Setting) {
	tmpl := template.Must(template.New("COMMON_TITLE").Parse(COMMON_TITLE))

	if err := tmpl.Execute(os.Stdout, conf); err != nil {
		panic(err)
		os.Exit(1)
	}

	confMap := structs.Map(conf)

	tmpl = template.Must(template.New("list_template").Parse(list_template))

	if err := tmpl.Execute(os.Stdout, confMap); err != nil {
		panic(err)
		os.Exit(1)
	}

}



