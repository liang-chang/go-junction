package config


//configs.Action 可用
var Actions = [...]string{"list", "check", "make", "recovery"}

//Symbolic.Action 可用
var SymbolicActions = [...]string{"ignore", "recovery"}

const FILE_SPLIT = "/"

type Setting struct {
	//list check create recovery
	Action    string

	Config    GlobalConfig
	PathAlias map[string]string
	Symbolic  []Symbolic
}

type GlobalConfig struct {
	//原始文件夹是否要重命名备份
	BackupLinkFolder   bool

	//是否要清空原始文件夹
	ClearBackupFolder  bool

	//#当target文件夹不存在时，创建文件夹或者终止
	CreateTargetFolder bool
}

type Symbolic struct {
	//类型，可选项: junction(只能针对文件夹) , hardlink(只能针对文件) , symbolic(两者都可以)
	//Type       string
	//ignore , recovery
	Action     string

	Target     string
	Link       []string

	LinkConfig []LinkConfig
}

//
type LinkConfig struct {
	FolderPattern   string

	MatchFolder      []string

	Backup          bool
	Clear           bool

	//即当最后一级文件目录不存在时，创建文件夹
	LastDirAppender bool

	//强制创建整个路径的文件夹
	ForeceCreate    bool
}