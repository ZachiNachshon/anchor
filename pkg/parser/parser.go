package parser

import "github.com/ZachiNachshon/anchor/pkg/models"

const (
	Identifier string = "parser"
)

type Parser interface {
	ParseInstructions(text string) (*models.InstructionsRoot, error)
}
