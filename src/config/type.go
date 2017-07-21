package config

//configs.Action 可用
var Actions = []string{"list", "check", "make", "recovery"}

//Symbolic.Action 可用
var SymbolicActions = []string{"ignore", "recovery"}

const FILE_SPLIT = "/"

type Setting struct {
	//list check create recovery
	Action string

	ConfigFile string

	Config    GlobalConfig
	PathAlias map[string]string
	Symbolic  []Symbolic
}

type GlobalConfig struct {
	//原始文件夹是否要重命名备份
	BackupLinkFolder bool

	//是否要清空原始文件夹
	ClearBackupFolder bool

	//#当target文件夹不存在时，创建文件夹或者终止
	CreateTargetFolder bool

	//target 文件夹路径表达式，与 filepath.Glob 用法一致
	targetFolders []string
}

type Symbolic struct {
	//类型，可选项: junction(只能针对文件夹) , hardlink(只能针对文件) , symbolic(两者都可以)
	//Type       string
	//ignore
	Skip bool

	Target string

	//初化化解析后，放入 []LinkConfig ,该不会再使用
	Link []string

	LinkConfig []LinkConfig
}

//
type LinkConfig struct {
	FolderPattern string

	MatchFolder []string

	Backup bool
	Clear  bool

	//是否在target文件夹中创建隔离文件夹
	Isolate bool

	//即当最后一级文件目录不存在时，创建文件夹
	LastDirAppender bool

	//强制创建整个路径的文件夹
	ForeceCreate bool

	//当目标文件夹不存在时，是否报警告
	WarnIgnore bool
}
