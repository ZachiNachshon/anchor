package printer

var CreateFakePrinterSpinner = func() *fakePrinterSpinner {
	return &fakePrinterSpinner{}
}

type fakePrinterSpinner struct {
	PrinterSpinner
	SpinMock          func()
	StopOnSuccessMock func()
	StopOnFailureMock func(err error)
}

func (pp *fakePrinterSpinner) Spin() {
	pp.SpinMock()
}

func (pp *fakePrinterSpinner) StopOnSuccess() {
	pp.StopOnSuccessMock()
}

func (pp *fakePrinterSpinner) StopOnFailure(err error) {
	pp.StopOnFailureMock(err)
}
