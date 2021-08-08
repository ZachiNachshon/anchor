package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/models"

	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"io/ioutil"
	"sort"
)

const (
	Identifier string = "extractor"
)

type Extractor interface {
	ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

type extractorImpl struct{}

func New() Extractor {
	return &extractorImpl{}
}

func (e *extractorImpl) ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
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
			sortInstructions(instructions.Instructions)
			sortWorkflows(instructions.Instructions)
			return instructions, nil
		}
	}
}

func sortInstructions(instructions *models.Instructions) {
	if instructions == nil || instructions.Actions == nil {
		return
	}
	sort.Slice(instructions.Actions, func(i, j int) bool {
		return instructions.Actions[i].Id < instructions.Actions[j].Id
	})
}

func sortWorkflows(instructions *models.Instructions) {
	if instructions == nil || instructions.Workflows == nil {
		return
	}
	sort.Slice(instructions.Workflows, func(i, j int) bool {
		return instructions.Workflows[i].Id < instructions.Workflows[j].Id
	})
}
