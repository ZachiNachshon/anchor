package parser

import "github.com/ZachiNachshon/anchor/pkg/models"

const (
	Identifier string = "parser"
)

type Parser interface {
	ParseCommandFolderInfo(text string) (*models.CommandFolderInfo, error)
	ParseInstructions(text string) (*models.InstructionsRoot, error)
}
