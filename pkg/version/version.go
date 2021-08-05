package version

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

// TODO: take from versioned file
const version = "v0.0.1"

func StartVersionVersionFlow(ctx common.Context) error {
	p, err := printer.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	p.PrintAnchorVersion(version)
	return nil
}
