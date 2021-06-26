package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

var configPrintTemplate = `Repository in-use: {{ .ConfigFilePath }}
---
{{ .ConfigText }}`

type PrintConfigTemplateItems struct {
	ConfigFilePath string
	ConfigText     string
}

type printerImpl struct {
	Printer
}

func New() Printer {
	return &printerImpl{}
}

func (p *printerImpl) PrintApplications(apps []*models.AppContent) {
	logger.Info("------ Applications ------")
	for _, app := range apps {
		logger.Infof("Name: %s", app)
	}
}

func (p *printerImpl) PrintConfiguration(ctx common.Context, cfgText string) {
	var items = PrintConfigTemplateItems{
		ConfigFilePath: ctx.AnchorFilesPath(),
		ConfigText:     cfgText,
	}
	if text, err := templates.TemplateToText(configPrintTemplate, items); err != nil {
		logger.Error("Failed to prepare configuration string")
	} else {
		fmt.Print(text)
	}
}
