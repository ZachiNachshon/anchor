package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

var configPrintTemplate = `---
Repository in-use: {{ .AnchorfilesRepoPath }}
---
Configuration file path: {{ .ConfigFilePath }} 
---
{{ .ConfigText }}`

type PrintConfigTemplateItems struct {
	AnchorfilesRepoPath string
	ConfigFilePath      string
	ConfigText          string
}

type printerImpl struct {
	Printer
}

func New() Printer {
	return &printerImpl{}
}

func (p *printerImpl) PrintApplications(apps []*models.ApplicationInfo) {
	logger.Info("------ Applications ------")
	for _, app := range apps {
		logger.Infof("Name: %s", app)
	}
}

func (p *printerImpl) PrintConfiguration(ctx common.Context, cfgFilePath string, cfgText string) {
	var items = PrintConfigTemplateItems{
		AnchorfilesRepoPath: ctx.AnchorFilesPath(),
		ConfigFilePath:      cfgFilePath,
		ConfigText:          cfgText,
	}
	if text, err := templates.TemplateToText(configPrintTemplate, items); err != nil {
		logger.Error("Failed to prepare configuration string")
	} else {
		fmt.Print(text)
	}
}
