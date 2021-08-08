package status

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

type AppStatusFunc func(ctx common.Context) error

var AppStatus = func(ctx common.Context) error {
	var l locator.Locator
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		l = resolved.(locator.Locator)
	}

	var p printer.Printer
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		p = resolved.(printer.Printer)
	}

	var e extractor.Extractor
	if resolved, err := ctx.Registry().SafeGet(extractor.Identifier); err != nil {
		return err
	} else {
		e = resolved.(extractor.Extractor)
	}

	var pa parser.Parser
	if resolved, err := ctx.Registry().SafeGet(parser.Identifier); err != nil {
		return err
	} else {
		pa = resolved.(parser.Parser)
	}

	p.PrintAnchorBanner()
	return runApplicationStatusFlow(ctx, l, e, pa, p)
}

func runApplicationStatusFlow(
	ctx common.Context,
	l locator.Locator,
	e extractor.Extractor,
	pa parser.Parser,
	p printer.Printer) error {

	var appStatus []*printer.AppStatusTemplateItem
	for _, app := range l.Applications() {
		status := &printer.AppStatusTemplateItem{
			Name: app.Name,
		}

		if !ioutils.IsValidPath(app.InstructionsPath) {
			status.MissingInstructionFile = true
		} else {
			inst, err := e.ExtractInstructions(app.InstructionsPath, pa)
			status.InvalidInstructionFormat = inst != nil || err != nil
		}

		status.CalculateValidity()
		appStatus = append(appStatus, status)
	}

	p.PrintApplications(appStatus)
	return nil
}
