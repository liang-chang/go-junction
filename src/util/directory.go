package util

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func DirectoryExist(path string) (bool, error) {
	attrs, err := getFileAttributes(path)
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	return attrs&syscall.FILE_ATTRIBUTE_DIRECTORY > 0, err
}

func IsSamePath(a, b string) bool {
	a = strings.ToLower(strings.Replace(a, `\`, `/`, -1))
	b = strings.ToLower(strings.Replace(b, `\`, `/`, -1))
	return strings.Compare(a, b) == 0
}

/**
dir 是否是 parent 的子文件或者是相同的文件夹
*/
func IsSubDirectory(dir, parent string) bool {
	dir = strings.ToLower(strings.Replace(dir, `\`, `/`, -1))
	parent = strings.ToLower(strings.Replace(parent, `\`, `/`, -1))
	ret := strings.Index(dir, parent) >= 0 && len(dir) != len(parent)
	return ret
}

func IsReparsePoint(path string) (bool, error) {
	attrs, err := getFileAttributes(path)
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	return attrs&syscall.FILE_ATTRIBUTE_REPARSE_POINT > 0, err
}

func Exist(path string) (bool, error) {
	attrs, err := getFileAttributes(path)
	if err != nil {
		return false, err
	}

	return (attrs&syscall.FILE_ATTRIBUTE_DIRECTORY > 0) || (attrs&syscall.FILE_ATTRIBUTE_DIRECTORY == 0), err
}

func FileExist(path string) (bool, error) {
	attrs, err := getFileAttributes(path)
	if err != nil {
		return false, err
	}

	return attrs&syscall.FILE_ATTRIBUTE_DIRECTORY == 0, err
}

//调用 windows 底层函数，判断文件路径是否是个文件夹
func getFileAttributes(path string) (attrs uint32, err error) {
	return syscall.GetFileAttributes(syscall.StringToUTF16Ptr(path))
}

func RemoveContents(dir string) error {
	var d *os.File
	var err error
	if d, err = os.Open(dir); err != nil {
		return err
	}
	defer d.Close()
	var names []string
	if names, err = d.Readdirnames(-1); err != nil {
		return err
	}
	for _, name := range names {
		if err = os.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}
	return nil
}

func getMatchDirectory(pattern string) (dirs []string, err error) {
	fils, err := filepath.Glob(pattern)

	if err != nil {
		return make([]string, 0, 0), err
	}

	ret := make([]string, 0, len(fils))
	for _, v := range fils {
		info, _ := os.Stat(v)
		if info.IsDir() {
			ret = append(ret, v)
		}
	}

	return ret, err
}
