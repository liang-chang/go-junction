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

type DoTarget func(target string, symbolic *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int)

type DoLink func(target string, linkIndex int, linkConfig *config.LinkConfig, symb *config.Symbolic, conf *config.Setting) (errCnt, warnCnt int)

func TraversalSymbolic(conf *config.Setting, symbolics  []config.Symbolic, doTarget DoTarget, doLink DoLink) (errCnt, warnCnt int) {
	errCnt = 0
	warnCnt = 0
	for sidex, symboT := range symbolics {

		var target string = symboT.Target
		currentSymbolic := &symbolics[sidex];

		e, w := doTarget(target, currentSymbolic, conf);
		errCnt += e
		warnCnt += w

		var linkConfigs []config.LinkConfig = currentSymbolic.LinkConfig

		for lindex, _ := range linkConfigs {

			linkConf := &linkConfigs[lindex]

			var matchFolder []string = linkConf.MatchFolder

			if len(matchFolder) == 0 {
				e, w := doLink(target, -1, linkConf, currentSymbolic, conf);
				errCnt += e
				warnCnt += w
				continue
			}

			for matchIndex, _ := range matchFolder {
				e, w := doLink(target, matchIndex, linkConf, currentSymbolic, conf);
				errCnt += e
				warnCnt += w
			}
		}
	}
	return errCnt, warnCnt
}