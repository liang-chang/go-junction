package config

type config struct {
	TargetConfig targetConfig
	PathAlias    map[string] string
	Junctions    []junction
}

type targetConfig struct {
	RenameTargetFolder bool
	ClearTargetFolder  bool
	SkipInvalidTarget  bool
}

type junction struct {
	Link   string
	Target string
}

