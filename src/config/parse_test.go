package config

import (
	"testing"
	//"fmt"
	//"unsafe"
	"strings"
	"fmt"
	"strconv"
)

func TestReadLinkText(t *testing.T) {
	var s = `bclf@d:/|\d+$|/bin`;
	index := strings.Index(s, "@")
	if index < 0 {
		index = 0
	}
	fmt.Println("%" + strconv.Itoa(10 + index) + "s %s")
}
