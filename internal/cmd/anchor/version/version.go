package version

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

// TODO: take from versioned file
const version = "v0.0.1"

type VersionVersionFunc func(ctx common.Context) error

var VersionVersion = func(ctx common.Context) error {
	var p printer.Printer
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		p = resolved.(printer.Printer)
		p.PrintAnchorVersion(version)
	}
	return nil
}
