package action

import (
//"util"
//"config"
//"os"
//"fmt"
////"junction"
////"syscall"
//"symbolic"
)
import (
	"config"
	"log"
	"os"
)

func Call(actionName string, conf config.Setting) {
	var fun action
	var ok bool
	if fun, ok = FUNC[actionName]; !ok {
		log.Fatal("unknown action: " + actionName)
		os.Exit(1)
	}
	fun(conf)
}
