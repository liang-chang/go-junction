package config

import (
	//"fmt"
	//"unsafe"
	//"strconv"
	//"strings"
	"fmt"
	"testing"
)

func TestReadLinkText(s *testing.T) {
	capitals := map[string]string{"France":"Paris", "Italy":"Rome", "Japan":"Tokyo" }
	for key := range capitals {
		fmt.Println("Map item: Capital of", key, "is", capitals[key])
	}
}