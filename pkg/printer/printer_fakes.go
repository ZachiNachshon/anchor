package printer

import (
	"github.com/ZachiNachshon/anchor/models"
)

var CreateFakePrinter = func() *fakePrinter {
	return &fakePrinter{}
}

type fakePrinter struct {
	Printer
	PrintAnchorBannerMock  func()
	PrintApplicationsMock  func(apps []*models.ApplicationInfo)
	PrintConfigurationMock func(cfgFilePath string, cfgText string)
}

func (p *fakePrinter) PrintAnchorBanner() {
	p.PrintAnchorBannerMock()
}

func (p *fakePrinter) PrintApplications(apps []*models.ApplicationInfo) {
	p.PrintApplicationsMock(apps)
}

func (p *fakePrinter) PrintConfiguration(cfgFilePath string, cfgText string) {
	p.PrintConfigurationMock(cfgFilePath, cfgText)
}
