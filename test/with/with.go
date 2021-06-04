package with

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/fakes"
	"os"
	"testing"
)

func Context(f func(ctx common.Context)) {
	ctx := common.EmptyAnchorContext()
	f(ctx)
}

func Logging(ctx common.Context, t *testing.T, f func(logger logger.Logger)) {
	createLogger(ctx, t, false, f)
}

func LoggingVerbose(ctx common.Context, t *testing.T, f func(logger logger.Logger)) {
	createLogger(ctx, t, true, f)
}

func createLogger(ctx common.Context, t *testing.T, verbose bool, f func(logger logger.Logger)) {
	if out, err := logger.FakeTestingLogger(ctx, t, verbose); err != nil {
		println("Failed to create a fake testing logger. error: %s", err)
		os.Exit(1)
	} else {
		logger.SetLogger(out)
		f(out)
	}
}

func Config(ctx common.Context, content string, f func(config config.AnchorConfig)) {
	if cfg, err := fakes.FakeConfigLoader(content); err != nil {
		logger.Fatalf("Failed to create a fake config loader. error: %s", err)
	} else {
		config.SetInContext(ctx, *cfg)
		f(*cfg)
	}
}
