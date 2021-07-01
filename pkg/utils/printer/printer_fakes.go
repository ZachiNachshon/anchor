package printer

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
)

var CreateFakePrinter = func() *fakePrinter {
	return &fakePrinter{}
}

type fakePrinter struct {
	Printer
	PrintApplicationsMock  func(apps []*models.AppContent)
	PrintConfigurationMock func(ctx common.Context, cfgFilePath string, cfgText string)
}

func (l *fakePrinter) PrintApplications(apps []*models.AppContent) {
	l.PrintApplicationsMock(apps)
}

func (l *fakePrinter) PrintConfiguration(ctx common.Context, cfgFilePath string, cfgText string) {
	l.PrintConfigurationMock(ctx, cfgFilePath, cfgText)
}
