package action

import (
	"config"
)

type action func(config.Setting)

var FUNC = map[string]action{
	"list":list,
	"check":check,
}