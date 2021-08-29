package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/models"

	"github.com/ZachiNachshon/anchor/pkg/parser"
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
	if contentByte, err := ioutil.ReadFile(instructionsPath); err != nil {
		return nil, fmt.Errorf("invalid instructions path. error: %s", err.Error())
	} else {
		var text = string(contentByte)

		if instRoot, err := p.ParseInstructions(text); err != nil {
			return nil, err
		} else {
			if instRoot != nil && instRoot.Instructions != nil {
				actions := instRoot.Instructions.Actions
				if actions != nil {
					sort.Slice(actions, func(i, j int) bool {
						return actions[i].Id < actions[j].Id
					})
				}

				workflows := instRoot.Instructions.Workflows
				if workflows != nil {
					sort.Slice(workflows, func(i, j int) bool {
						return workflows[i].Id < workflows[j].Id
					})
				}
			}
			return instRoot, nil
		}
	}
}
