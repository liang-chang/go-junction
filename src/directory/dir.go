package directory

import (
	"os"
	//"os/exec"
	//"bytes"
	//"fmt"
	"fmt"
)

func DirectoryExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	fmt.Println(fileInfo)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func cmdExec(cmdStr string, args ...string) {
	//cmd := exec.Command(cmdStr,args)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Run()
	//return out.String()
}
