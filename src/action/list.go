package action

import (
	"fmt"
	"config"
)

func list(conf config.Setting) {
	for _, v := range conf.Symbolic {
		fmt.Println(v.LinkConfig)
	}
}

