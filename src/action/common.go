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

type DoTarget func(target string, symbolic *config.Symbolic) (errCnt,warnCnt int)

type DoLink func(target, link string, folderIndex int, linkConfig *config.LinkConfig) (errCnt,warnCnt int)

func TraversalSymbolic(symbolics  []config.Symbolic, doTarget DoTarget, doLink DoLink) (errCnt, warnCnt int) {
	errCnt = 0
	warnCnt = 0
	for sidex, symboT := range symbolics {

		var target string = symboT.Target

		errCnt += doTarget(target, &symbolics[sidex]);

		var linkConfigs []config.LinkConfig = symbolics[sidex].LinkConfig

		for lindex, _ := range linkConfigs {

			linkConf := &linkConfigs[lindex]

			var matchFolder []string = linkConf.MatchFolder

			if len(matchFolder) == 0 {
				warnCnt++
				linkConf.MatchFolder = append(linkConf.MatchFolder, "Warning! No directory match !")
				continue
			}

			for matchIndex, link := range matchFolder {
				errCnt += doLink(target, link, matchIndex, linkConf);
			}
		}
	}
	return errCnt, warnCnt
}