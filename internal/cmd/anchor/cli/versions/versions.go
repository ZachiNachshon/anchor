package versions

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

type CliVersionsFunc func(ctx common.Context, o *versionsOrchestrator) error

var CliVersions = func(ctx common.Context, o *versionsOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	o.bannerFunc(o)
	return o.runFunc(o, ctx)
}

type versionsOrchestrator struct {
	prntr printer.Printer

	prepareFunc func(o *versionsOrchestrator, ctx common.Context) error
	bannerFunc  func(o *versionsOrchestrator)
	runFunc     func(o *versionsOrchestrator, ctx common.Context) error
}

func NewOrchestrator() *versionsOrchestrator {
	return &versionsOrchestrator{
		bannerFunc:  banner,
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *versionsOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func banner(o *versionsOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *versionsOrchestrator, ctx common.Context) error {
	return nil
}
