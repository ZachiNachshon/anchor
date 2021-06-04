package extractor

import "github.com/ZachiNachshon/anchor/pkg/utils/parser"

type Extractor interface {
	ExtractPromptItems(instructionsPath string) (*parser.PromptItems, error)
}
