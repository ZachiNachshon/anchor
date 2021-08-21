package install

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

type ControllerInstallFunc func(ctx common.Context, o *installOrchestrator) error

var ControllerInstall = func(ctx common.Context, o *installOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	o.bannerFunc(o)
	return o.runFunc(o, ctx)
}

type installOrchestrator struct {
	prntr printer.Printer

	prepareFunc func(o *installOrchestrator, ctx common.Context) error
	bannerFunc  func(o *installOrchestrator)
	runFunc     func(o *installOrchestrator, ctx common.Context) error
}

func NewOrchestrator() *installOrchestrator {
	return &installOrchestrator{
		bannerFunc:  banner,
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *installOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func banner(o *installOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *installOrchestrator, ctx common.Context) error {
	return nil
}
