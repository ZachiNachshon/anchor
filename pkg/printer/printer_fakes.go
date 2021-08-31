package printer

var CreateFakePrinter = func() *fakePrinter {
	return &fakePrinter{}
}

type fakePrinter struct {
	Printer
	PrintAnchorBannerMock       func()
	PrintAnchorVersionMock      func(version string)
	PrintApplicationsStatusMock func(apps []*AppStatusTemplateItem)
	PrintConfigurationMock      func(cfgFilePath string, cfgText string)
	PrepareRunActionPlainerMock func(actionId string) PrinterPlainer
	PrepareRunActionSpinnerMock func(actionId string, scriptOutputPath string) PrinterSpinner
}

func (p *fakePrinter) PrintAnchorBanner() {
	p.PrintAnchorBannerMock()
}

func (p *fakePrinter) PrintAnchorVersion(version string) {
	p.PrintAnchorVersionMock(version)
}

func (p *fakePrinter) PrintApplicationsStatus(apps []*AppStatusTemplateItem) {
	p.PrintApplicationsStatusMock(apps)
}

func (p *fakePrinter) PrintConfiguration(cfgFilePath string, cfgText string) {
	p.PrintConfigurationMock(cfgFilePath, cfgText)
}

func (p *fakePrinter) PrepareRunActionPlainer(actionId string) PrinterPlainer {
	return p.PrepareRunActionPlainerMock(actionId)
}

func (p *fakePrinter) PrepareRunActionSpinner(actionId string, scriptOutputPath string) PrinterSpinner {
	return p.PrepareRunActionSpinnerMock(actionId, scriptOutputPath)
}
