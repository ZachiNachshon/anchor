package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

func StartApplicationStatusFlow(ctx common.Context) error {
	l, err := locator.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}

	p, err := printer.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}

	e, err := extractor.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}

	pa, err := parser.FromRegistry(ctx.Registry())
	if err != nil {
		return err
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
