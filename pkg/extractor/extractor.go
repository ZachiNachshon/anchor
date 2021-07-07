package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"io/ioutil"
	"sort"
)

type extractorImpl struct{}

func New() Extractor {
	return &extractorImpl{}
}

func (e *extractorImpl) ExtractPromptItems(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
	if !ioutils.IsValidPath(instructionsPath) {
		return nil, fmt.Errorf("invalid instructions path. path: %s", instructionsPath)
	}

	if contentByte, err := ioutil.ReadFile(instructionsPath); err != nil {
		return nil, err
	} else {
		var text = string(contentByte)

		if instructions, err := p.ParseInstructions(text); err != nil {
			return nil, err
		} else {
			return sortInstructions(instructions), nil
		}
	}
}

func sortInstructions(instructions *models.Instructions) *models.Instructions {
	sort.Slice(instructions.Items, func(i, j int) bool {
		return instructions.Items[i].Id < instructions.Items[j].Id
	})
	return instructions
}
