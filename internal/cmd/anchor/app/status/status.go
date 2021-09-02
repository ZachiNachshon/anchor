package status

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

type AppStatusFunc func(ctx common.Context, o *statusOrchestrator) error

var AppStatus = func(ctx common.Context, o *statusOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	o.bannerFunc(o)
	return o.runFunc(o, ctx)
}

type statusOrchestrator struct {
	l                 locator.Locator
	e                 extractor.Extractor
	prsr              parser.Parser
	prntr             printer.Printer
	validStatusOnly   bool
	invalidStatusOnly bool

	// --- CLI Command ---
	prepareFunc func(o *statusOrchestrator, ctx common.Context) error
	bannerFunc  func(o *statusOrchestrator)
	runFunc     func(o *statusOrchestrator, ctx common.Context) error
}

func NewOrchestrator() *statusOrchestrator {
	return &statusOrchestrator{
		// --- CLI Command ---
		bannerFunc:  banner,
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *statusOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		o.l = resolved.(locator.Locator)
	}

	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}

	if resolved, err := ctx.Registry().SafeGet(extractor.Identifier); err != nil {
		return err
	} else {
		o.e = resolved.(extractor.Extractor)
	}

	if resolved, err := ctx.Registry().SafeGet(parser.Identifier); err != nil {
		return err
	} else {
		o.prsr = resolved.(parser.Parser)
	}

	return nil
}

func banner(o *statusOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *statusOrchestrator, ctx common.Context) error {
	var appStatus []*printer.AppStatusTemplateItem
	for _, app := range o.l.Applications() {
		status := &printer.AppStatusTemplateItem{
			Name: app.Name,
		}

		if !ioutils.IsValidPath(app.InstructionsPath) {
			status.MissingInstructionFile = true
		} else {
			inst, err := o.e.ExtractInstructions(app.InstructionsPath, o.prsr)
			status.InvalidInstructionFormat = inst == nil || err != nil
		}

		isValid := status.CalculateValidity()
		if isValid && o.validStatusOnly ||
			!isValid && o.invalidStatusOnly ||
			!o.validStatusOnly && !o.invalidStatusOnly {

			appStatus = append(appStatus, status)
		}
	}

	o.prntr.PrintApplicationsStatus(appStatus)
	return nil
}
