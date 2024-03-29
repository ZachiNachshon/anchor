package version

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/resources"
)

type VersionVersionFunc func(ctx common.Context, o *versionOrchestrator) error

var VersionVersion = func(ctx common.Context, o *versionOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	return o.runFunc(o, ctx)
}

type versionOrchestrator struct {
	prntr printer.Printer

	prepareFunc func(o *versionOrchestrator, ctx common.Context) error
	runFunc     func(o *versionOrchestrator, ctx common.Context) error
}

func NewOrchestrator() *versionOrchestrator {
	return &versionOrchestrator{
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *versionOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func run(o *versionOrchestrator, ctx common.Context) error {
	if version, err := resources.VersionFile.ReadFile("version.txt"); err != nil {
		return err
	} else {
		o.prntr.PrintAnchorVersion(string(version))
		return nil
	}
}
