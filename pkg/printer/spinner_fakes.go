package printer

var CreateFakePrinterSpinner = func() *fakePrinterSpinner {
	return &fakePrinterSpinner{}
}

type fakePrinterSpinner struct {
	PrinterSpinner
	SpinMock                           func()
	StopOnSuccessMock                  func()
	StopOnSuccessWithCustomMessageMock func(message string)
	StopOnFailureMock                  func(err error)
}

func (pp *fakePrinterSpinner) Spin() {
	pp.SpinMock()
}

func (pp *fakePrinterSpinner) StopOnSuccess() {
	pp.StopOnSuccessMock()
}

func (pp *fakePrinterSpinner) StopOnSuccessWithCustomMessage(message string) {
	pp.StopOnSuccessWithCustomMessageMock(message)
}

func (pp *fakePrinterSpinner) StopOnFailure(err error) {
	pp.StopOnFailureMock(err)
}
