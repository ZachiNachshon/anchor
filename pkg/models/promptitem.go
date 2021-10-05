package models

type Action struct {
	Id                  string `yaml:"id"`
	Title               string `yaml:"title"`
	Description         string `yaml:"description"`
	Script              string `yaml:"script"`
	ScriptFile          string `yaml:"scriptFile"`
	ShowOutput          bool   `yaml:"showOutput"`
	Context             string `yaml:"context"`
	AnchorfilesRepoPath string `yaml:"-"` // Used as a working directory for script file execution
}

type Workflow struct {
	Id               string   `yaml:"id"`
	Title            string   `yaml:"title"`
	Description      string   `yaml:"description"`
	TolerateFailures bool     `yaml:"tolerateFailures"`
	Context          string   `yaml:"context"`
	ActionIds        []string `yaml:"actionIds"`
}
