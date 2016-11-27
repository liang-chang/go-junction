package action

import (
	"config"
)

type action func(config.Setting)

var FUNC = map[string]action{
	"list":list,
	"check":check,
	"make":make,
}

const FOLDER_BACK_SUBFFIX = "_junction_bak_$_$"