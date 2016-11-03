package util

import (
	"syscall"
	//"fmt"
)

/**
调用 windows 底层函数，判断文件路径是否是个文件夹
 */
func DirectoryExist(path string) (bool, error) {
	attrs, err := syscall.GetFileAttributes(syscall.StringToUTF16Ptr(path));
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	return attrs & syscall.FILE_ATTRIBUTE_DIRECTORY > 0, nil
}

func cmdExec(cmdStr string, args ...string) {
	//cmd := exec.Command(cmdStr,args)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Run()
	//return out.String()
}
