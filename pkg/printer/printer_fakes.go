package printer

var CreateFakePrinter = func() *fakePrinter {
	return &fakePrinter{}
}

type fakePrinter struct {
	Printer
	PrintAnchorBannerMock                      func()
	PrintAnchorVersionMock                     func(version string)
	PrintAnchorFolderItemStatusMock            func(anchorFolderItemsStatus []*AnchorFolderItemStatusTemplate)
	PrintConfigurationMock                     func(cfgFilePath string, cfgText string)
	PrintMissingInstructionsMock               func()
	PrintEmptyLinesMock                        func(count int)
	PrintSuccessMock                           func(message string)
	PrintWarningMock                           func(message string)
	PrintErrorMock                             func(message string)
	PrepareRunActionPlainerMock                func(actionId string) PrinterPlainer
	PrepareRunActionSpinnerMock                func(actionId string, scriptOutputPath string) PrinterSpinner
	PrepareReadRemoteHeadCommitHashSpinnerMock func(url string, branch string) PrinterSpinner
	PrepareCloneRepositorySpinnerMock          func(url string, branch string) PrinterSpinner
	PrepareResetToRevisionSpinnerMock          func(revision string) PrinterSpinner
}

func (p *fakePrinter) PrintAnchorBanner() {
	p.PrintAnchorBannerMock()
}

func (p *fakePrinter) PrintAnchorVersion(version string) {
	p.PrintAnchorVersionMock(version)
}

func (p *fakePrinter) PrintAnchorFolderItemStatus(apps []*AnchorFolderItemStatusTemplate) {
	p.PrintAnchorFolderItemStatusMock(apps)
}

func (p *fakePrinter) PrintConfiguration(cfgFilePath string, cfgText string) {
	p.PrintConfigurationMock(cfgFilePath, cfgText)
}

func (p *fakePrinter) PrintMissingInstructions() {
	p.PrintMissingInstructionsMock()
}

func (p *fakePrinter) PrintEmptyLines(count int) {
	p.PrintEmptyLinesMock(count)
}

func (p *fakePrinter) PrintSuccess(message string) {
	p.PrintSuccessMock(message)
}

func (p *fakePrinter) PrintWarning(message string) {
	p.PrintWarningMock(message)
}

func (p *fakePrinter) PrintError(message string) {
	p.PrintErrorMock(message)
}

func (p *fakePrinter) PrepareRunActionPlainer(actionId string) PrinterPlainer {
	return p.PrepareRunActionPlainerMock(actionId)
}

func (p *fakePrinter) PrepareRunActionSpinner(actionId string, scriptOutputPath string) PrinterSpinner {
	return p.PrepareRunActionSpinnerMock(actionId, scriptOutputPath)
}

func (p *fakePrinter) PrepareReadRemoteHeadCommitHashSpinner(url string, branch string) PrinterSpinner {
	return p.PrepareReadRemoteHeadCommitHashSpinnerMock(url, branch)
}

func (p *fakePrinter) PrepareCloneRepositorySpinner(url string, branch string) PrinterSpinner {
	return p.PrepareCloneRepositorySpinnerMock(url, branch)
}

func (p *fakePrinter) PrepareResetToRevisionSpinner(revision string) PrinterSpinner {
	return p.PrepareResetToRevisionSpinnerMock(revision)
}
