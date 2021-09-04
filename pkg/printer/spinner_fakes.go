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
	StopOnFailureWithCustomMessageMock func(message string)
}

func (ps *fakePrinterSpinner) Spin() {
	ps.SpinMock()
}

func (ps *fakePrinterSpinner) StopOnSuccess() {
	ps.StopOnSuccessMock()
}

func (ps *fakePrinterSpinner) StopOnSuccessWithCustomMessage(message string) {
	ps.StopOnSuccessWithCustomMessageMock(message)
}

func (ps *fakePrinterSpinner) StopOnFailure(err error) {
	ps.StopOnFailureMock(err)
}

func (ps *fakePrinterSpinner) StopOnFailureWithCustomMessage(message string) {
	ps.StopOnFailureWithCustomMessageMock(message)
}
