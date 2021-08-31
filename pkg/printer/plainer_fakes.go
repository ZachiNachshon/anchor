package printer

var CreateFakePrinterPlainer = func() *fakePrinterPlainer {
	return &fakePrinterPlainer{}
}

type fakePrinterPlainer struct {
	PrinterPlainer
	StartMock         func()
	StopOnSuccessMock func()
	StopOnFailureMock func(err error)
}

func (pp *fakePrinterPlainer) Start() {
	pp.StartMock()
}

func (pp *fakePrinterPlainer) StopOnSuccess() {
	pp.StopOnSuccessMock()
}

func (pp *fakePrinterPlainer) StopOnFailure(err error) {
	pp.StopOnFailureMock(err)
}
