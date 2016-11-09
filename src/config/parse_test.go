package config

import (
	"testing"
	"fmt"
	"regexp"
	"strings"
)

func TestReadLinkText(t *testing.T) {

	diskReg := regexp.MustCompile(`{{[^}]+}}`)

	parts := diskReg.FindAllStringSubmatch(`V:\{{chrome}}\Defaultgfasdgfasdgasdgfasdgas{{temp}}`, -1)
	for _, v := range parts {
		fmt.Println(strings.TrimRight(strings.TrimLeft(v[0],"{{"),"}}"))
	}

	//readLinkText(`bclf@d:/|\d+$|/bin`)
}