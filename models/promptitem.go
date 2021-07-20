package models

type Action struct {
	Id          string `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	File        string `yaml:"file"`
}

type Workflow struct {
	Id               string   `yaml:"id"`
	Description      string   `yaml:"description"`
	TolerateFailures bool     `yaml:"tolerateFailures"`
	ActionIds        []string `yaml:"actionIds"`
}
