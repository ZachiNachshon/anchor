package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"io/ioutil"
	"os"
	"sort"
)

const (
	Identifier string = "extractor"
)

type Extractor interface {
	ExtractCommandFolderInfo(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error)
	ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

type extractorImpl struct{}

func New() Extractor {
	return &extractorImpl{}
}

func (e *extractorImpl) ExtractCommandFolderInfo(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error) {
	commandFilePath := fmt.Sprintf("%s/%s", dirPath, globals.AnchorCommandFileName)
	if contentByte, err := ioutil.ReadFile(commandFilePath); err != nil {
		return nil, fmt.Errorf("invalid anchor folder info path. error: %s", err.Error())
	} else {
		var text = string(contentByte)
		// expand potential environment variables
		contentByteExpanded := []byte(os.ExpandEnv(text))

		if commandFolder, err := p.ParseCommandFolderInfo(string(contentByteExpanded)); err != nil {
			return nil, err
		} else {
			if commandFolder != nil &&
				len(commandFolder.Command.Use) > 0 &&
				len(commandFolder.Command.Short) > 0 {
				commandFolder.DirPath = dirPath
				return commandFolder, nil
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
		// expand potential environment variables
		contentByteExpanded := []byte(os.ExpandEnv(text))

		if instRoot, err := p.ParseInstructions(string(contentByteExpanded)); err != nil {
			return nil, err
		} else {
			if instRoot != nil && instRoot.Instructions != nil {
				actions := instRoot.Instructions.Actions
				if actions != nil {
					sort.Slice(actions, func(i, j int) bool {
						if len(actions[i].DisplayName) > 0 && len(actions[j].DisplayName) > 0 {
							return actions[i].DisplayName < actions[j].DisplayName
						} else if len(actions[i].DisplayName) > 0 && len(actions[j].DisplayName) == 0 {
							return actions[i].DisplayName < actions[j].Id
						} else if len(actions[i].DisplayName) == 0 && len(actions[j].DisplayName) > 0 {
							return actions[i].Id < actions[j].DisplayName
						}
						return actions[i].Id < actions[j].Id
					})
				}

				workflows := instRoot.Instructions.Workflows
				if workflows != nil {
					sort.Slice(workflows, func(i, j int) bool {
						if len(workflows[i].DisplayName) > 0 && len(workflows[j].DisplayName) > 0 {
							return workflows[i].DisplayName < workflows[j].DisplayName
						} else if len(workflows[i].DisplayName) > 0 && len(workflows[j].DisplayName) == 0 {
							return workflows[i].DisplayName < workflows[j].Id
						} else if len(workflows[i].DisplayName) == 0 && len(workflows[j].DisplayName) > 0 {
							return workflows[i].Id < workflows[j].DisplayName
						}
						return workflows[i].Id < workflows[j].Id
					})
				}
			}
			return instRoot, nil
		}
	}
}
