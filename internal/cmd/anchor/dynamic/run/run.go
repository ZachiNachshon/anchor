package run

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

type DynamicRunFunc func(ctx common.Context, o *runOrchestrator) error

var DynamicRun = func(ctx common.Context, o *runOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	o.bannerFunc(o)
	return o.runFunc(o, ctx)
}

type runOrchestrator struct {
	parentFolderName string

	prntr printer.Printer

	prepareFunc func(o *runOrchestrator, ctx common.Context) error
	bannerFunc  func(o *runOrchestrator)
	runFunc     func(o *runOrchestrator, ctx common.Context) error
}

func NewOrchestrator(parentFolderName string) *runOrchestrator {
	return &runOrchestrator{
		parentFolderName: parentFolderName,
		bannerFunc:       banner,
		prepareFunc:      prepare,
		runFunc:          run,
	}
}

func prepare(o *runOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func banner(o *runOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *runOrchestrator, ctx common.Context) error {
	return nil
}
