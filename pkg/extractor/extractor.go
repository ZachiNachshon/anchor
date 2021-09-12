package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"io/ioutil"
	"sort"
)

const (
	Identifier string = "extractor"
)

type Extractor interface {
	ExtractAnchorFolderInfo(dirPath string, p parser.Parser) (*models.AnchorFolderInfo, error)
	ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

type extractorImpl struct{}

func New() Extractor {
	return &extractorImpl{}
}

func (e *extractorImpl) ExtractAnchorFolderInfo(dirPath string, p parser.Parser) (*models.AnchorFolderInfo, error) {
	commandFilePath := fmt.Sprintf("%s/%s", dirPath, globals.AnchorCommandFileName)
	if contentByte, err := ioutil.ReadFile(commandFilePath); err != nil {
		return nil, fmt.Errorf("invalid anchor folder info path. error: %s", err.Error())
	} else {
		var text = string(contentByte)

		if anchorFolder, err := p.ParseAnchorFolderInfo(text); err != nil {
			return nil, err
		} else {
			if anchorFolder != nil &&
				len(anchorFolder.Command.Use) > 0 &&
				len(anchorFolder.Command.Short) > 0 {
				anchorFolder.DirPath = dirPath
				return anchorFolder, nil
			}
			return nil, fmt.Errorf("bad anchor folder command file structure")
		}
	}
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
