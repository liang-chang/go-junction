package config

import (
	"testing"
	"fmt"
)

func TestReadLinkText(t *testing.T) {
	var conf Symbolic
	linkConfig := make([]LinkConfig, 0, 5);
	linkConfig = append(linkConfig, LinkConfig{})
	fmt.Println(conf)
	//readLinkText(`bclf@d:/|\d+$|/bin`)
}