package util

import (
	"syscall"
	//"fmt"
)

func DirectoryExist(path string) (bool, error) {
	attrs, err := getFileAttributes(path);
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	return attrs & syscall.FILE_ATTRIBUTE_DIRECTORY > 0, err
}

func FileExist(path string) (bool, error) {
	attrs, err := getFileAttributes(path);
	if err != nil {
		return false, err
	}

	return attrs & syscall.FILE_ATTRIBUTE_DIRECTORY == 0, err
}


//调用 windows 底层函数，判断文件路径是否是个文件夹
func getFileAttributes(path string) (attrs uint32, err error) {
	return syscall.GetFileAttributes(syscall.StringToUTF16Ptr(path));
}

func cmdExec(cmdStr string, args ...string) {
	//cmd := exec.Command(cmdStr,args)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Run()
	//return out.String()
}
