package version

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VersionActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "print version successfully",
			Func: PrintVersionSuccessfully,
		},
		{
			Name: "fail to print version",
			Func: FailToPrintVersion,
		},
	}
	harness.RunTests(t, tests)
}

var PrintVersionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintAnchorVersionMock = func(version string) {}
			ctx.Registry().Set(printer.Identifier, fakePrinter)
			err := VersionVersion(ctx)
			assert.Nil(t, err, "expected print version to succeed")
		})
	})
}

var FailToPrintVersion = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := VersionVersion(ctx)
			assert.NotNil(t, err, "expected print version to fail")
		})
	})
}
