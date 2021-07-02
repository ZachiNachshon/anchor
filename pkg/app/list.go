package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

func StartApplicationListFlow(ctx common.Context) error {
	l, err := locator.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}

	p, err := printer.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}

	p.PrintApplications(l.Applications())
	return nil
}
