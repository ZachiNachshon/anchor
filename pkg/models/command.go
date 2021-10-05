package models

type AnchorFolderInfo struct {
	Name        string               `yaml:"name"`
	Command     *AnchorFolderCommand `yaml:"command"`
	Description string               `yaml:"description"`
	Items       map[string]*AnchorFolderItemInfo
	DirPath     string
}

type AnchorFolderCommand struct {
	Use   string `yaml:"use"`
	Short string `yaml:"short"`
}

type AnchorFolderItemInfo struct {
	Name             string
	DirPath          string
	InstructionsPath string
}
