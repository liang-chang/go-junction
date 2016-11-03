package config

type config struct {
	//list create recovery
	ActionName string

	LinkConfig linkConfig
	PathAlias  map[string]string
	Junction   []symbolic
}

type linkConfig struct {
	// 原始文件夹是否要重命名备份
	RenameTargetFolder bool

	//是否要清空原始文件夹
	ClearTargetFolder  bool

	//#当目录文件夹无效时是否继续
	SkipInvalidTarget  bool
}

type symbolic struct {
	//类型，可选项: junction(只能针对文件夹) , hardlink(只能针对文件) , symbolic(两者都可以)
	Type   string

	//create , ignore , recovery
	Action bool
	Target string
	Link   []string
}

