package config

import (
	"testing"
	"strings"
	"fmt"
)

func TestReadLinkText(t *testing.T) {
	split := strings.Split(`bclf@d:/|\d+$|/bin`, "@")
	fmt.Println(split[:])
	fmt.Println(split[1:])
	//readLinkText(`bclf@d:/|\d+$|/bin`)
}