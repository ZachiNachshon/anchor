package status

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

type AppStatusFunc func(ctx common.Context, o orchestrator) error

var AppStatus = func(ctx common.Context, o orchestrator) error {
	err := o.prepare(ctx)
	if err != nil {
		return err
	}
	o.banner()
	return o.run(ctx)
}

type orchestrator interface {
	banner()
	prepare(ctx common.Context) error
	run(ctx common.Context) error
}

type statusOrchestratorImpl struct {
	orchestrator
	l  locator.Locator
	e  extractor.Extractor
	pa parser.Parser
	p  printer.Printer
}

var statusOrchestrator = &statusOrchestratorImpl{}

func (o *statusOrchestratorImpl) banner() {
	o.p.PrintAnchorBanner()
}

func (o *statusOrchestratorImpl) prepare(ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		o.l = resolved.(locator.Locator)
	}

	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.p = resolved.(printer.Printer)
	}

	if resolved, err := ctx.Registry().SafeGet(extractor.Identifier); err != nil {
		return err
	} else {
		o.e = resolved.(extractor.Extractor)
	}

	if resolved, err := ctx.Registry().SafeGet(parser.Identifier); err != nil {
		return err
	} else {
		o.pa = resolved.(parser.Parser)
	}

	return nil
}

func (o *statusOrchestratorImpl) run(ctx common.Context) error {
	var appStatus []*printer.AppStatusTemplateItem
	for _, app := range o.l.Applications() {
		status := &printer.AppStatusTemplateItem{
			Name: app.Name,
		}

		if !ioutils.IsValidPath(app.InstructionsPath) {
			status.MissingInstructionFile = true
		} else {
			inst, err := o.e.ExtractInstructions(app.InstructionsPath, o.pa)
			status.InvalidInstructionFormat = inst == nil || err != nil
		}

		status.CalculateValidity()
		appStatus = append(appStatus, status)
	}

	o.p.PrintApplications(appStatus)
	return nil
}
