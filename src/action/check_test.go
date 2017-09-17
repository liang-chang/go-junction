package action

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetMatchDirectory(t *testing.T) {
	fmt.Print(filepath.Dir(os.Args[0]))

}
