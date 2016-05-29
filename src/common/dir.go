package common

import (
	"os"
	//"os/exec"
	//"bytes"
	//"fmt"
)

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func cmdExec(cmdStr string,args...string){
	//cmd := exec.Command(cmdStr,args)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Run()
	//return out.String()
}
