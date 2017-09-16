package action

import (
	"os"
	"testing"
	"fmt"
	"path/filepath"
)

func TestGetMatchDirectory(t *testing.T) {
	fmt.Print(filepath.Dir(os.Args[0]))

}