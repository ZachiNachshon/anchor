package printer

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"time"
)

type PrinterSpinner interface {
	Spin()
	StopOnSuccess()
	StopOnFailure(err error)
}

type printerSpinnerImpl struct {
	PrinterSpinner
	spnr *spinner.Spinner

	runMsg           string
	successMsg       string
	failureMsgFormat string
}

func NewSpinner(runMsg string, successMsg string, failureMsgFormat string) PrinterSpinner {
	spnr := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	_ = spnr.Color("blue")
	spnr.HideCursor = true
	spnr.Prefix = ""
	return &printerSpinnerImpl{
		spnr:             spnr,
		runMsg:           runMsg,
		successMsg:       successMsg,
		failureMsgFormat: failureMsgFormat,
	}
}

func (p *printerSpinnerImpl) Spin() {
	p.spnr.Suffix = p.runMsg
	p.spnr.Start()
}

func (p *printerSpinnerImpl) StopOnSuccess() {
	_, _ = fmt.Fprintf(os.Stdout, "\r \r")
	fmt.Printf(p.successMsg)
	p.spnr.Stop()
}

func (p *printerSpinnerImpl) StopOnFailure(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "\r \r")
	fmt.Printf(fmt.Sprintf(p.failureMsgFormat, err.Error()))
	p.spnr.Stop()
}
