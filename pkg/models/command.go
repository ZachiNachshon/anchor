package models

type CommandFolderInfo struct {
	Name        string                `yaml:"name"`
	Command     *CommandFolderCommand `yaml:"command"`
	Description string                `yaml:"description"`
	Items       map[string]*CommandFolderItemInfo
	DirPath     string
}

type CommandFolderCommand struct {
	Use   string `yaml:"use"`
	Short string `yaml:"short"`
}

type CommandFolderItemInfo struct {
	Name             string
	DirPath          string
	InstructionsPath string
}
