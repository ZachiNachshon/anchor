package models

type AnchorFolderInfo struct {
	Type    string
	Name    string
	Command *AnchorFolderCommand
	DirPath string
	Items   map[string]*AnchorFolderItemInfo
}

type AnchorFolderCommand struct {
	Use   string
	Short string
}

type AnchorFolderItemInfo struct {
	Name             string
	DirPath          string
	InstructionsPath string
}
