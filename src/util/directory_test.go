package util

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestGetMatchDirectory(t *testing.T) {
	var baseDir = `v:/useless/`
	var size = 10
	for i := 0; i < size; i++ {
		os.MkdirAll(baseDir+`Z/`+strconv.Itoa(i), os.ModePerm)
	}

	defer func() {
		for i := 0; i < size; i++ {
			os.RemoveAll(baseDir + `Z/` + strconv.Itoa(i))
		}
	}()

	dirs, err := getMatchDirectory(baseDir + "?/?")
	if err != nil {
		t.Errorf("出错了")
	}
	if len(dirs) < size {
		t.Errorf("取出数量有误")
	}
	Log(dirs)
}

func TestFilePathGlob(t *testing.T) {
	matches, err := filepath.Glob(`v:/useless/[A-Z]/Z[0-9]`)
	Log(err)
	Log(matches)
}
