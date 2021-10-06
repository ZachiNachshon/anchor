package printer

import (
	"fmt"
)

type PrinterPlainer interface {
	Start()
	StopOnSuccess()
	StopOnFailure(err error)
}

type printerPlainerImpl struct {
	PrinterPlainer

	runMsg           string
	successMsg       string
	failureMsgFormat string
}

func NewPlainer(runMsg string, successMsg string, failureMsg string) PrinterPlainer {
	return &printerPlainerImpl{
		runMsg:           runMsg,
		successMsg:       successMsg,
		failureMsgFormat: failureMsg,
	}
}

func (p *printerPlainerImpl) Start() {
	fmt.Println(p.runMsg)
}

func (p *printerPlainerImpl) StopOnSuccess() {
	//_, _ = fmt.Fprintf(os.Stdout, "\r \r")
	fmt.Printf(p.successMsg)
}

func (p *printerPlainerImpl) StopOnFailure(err error) {
	//_, _ = fmt.Fprintf(os.Stdout, "\r \r")
	fmt.Printf(fmt.Sprintf(p.failureMsgFormat, err.Error()))
}
