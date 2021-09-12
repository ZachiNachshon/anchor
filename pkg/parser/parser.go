package parser

import "github.com/ZachiNachshon/anchor/pkg/models"

const (
	Identifier string = "parser"
)

type Parser interface {
	ParseAnchorFolderInfo(text string) (*models.AnchorFolderInfo, error)
	ParseInstructions(text string) (*models.InstructionsRoot, error)
}
