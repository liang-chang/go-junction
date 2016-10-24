package config

type config struct {
	ActionName   string
	TargetConfig targetConfig
	PathAlias    map[string]string
	Junction     []junction
}

type targetConfig struct {
	// 原始文件夹是否要重命名备份
	RenameTargetFolder bool

	//是否要清空原始文件夹
	ClearTargetFolder  bool

	//#当目录文件夹无效时是否继续
	SkipInvalidTarget  bool
}

type junction struct {
	Target string
	Link   []string
}

