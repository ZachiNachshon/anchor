package with

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"os"
	"testing"
)

func Context(f func(ctx common.Context)) {
	// We want a new registry instance on every context
	ctx := common.EmptyAnchorContext(registry.New())
	f(ctx)
}

func Logging(ctx common.Context, t *testing.T, f func(logger logger.Logger)) {
	createLogger(ctx, t, false, f)
}

func LoggingVerbose(ctx common.Context, t *testing.T, f func(logger logger.Logger)) {
	createLogger(ctx, t, true, f)
}

func createLogger(ctx common.Context, t *testing.T, verbose bool, f func(logger logger.Logger)) {
	if out, err := logger.FakeTestingLogger(t, verbose); err != nil {
		println("Failed to create a fake testing logger. error: %s", err)
		os.Exit(1)
	} else {
		logger.SetLogger(out)
		f(out)
	}
}

func Config(ctx common.Context, content string, f func(config config.AnchorConfig)) {
	if cfg, err := config.ViperConfigInMemoryLoader(content); err != nil {
		logger.Fatalf("Failed to create a fake config loader. error: %s", err)
	} else {
		config.SetInContext(ctx, *cfg)
		// set current config context as the active config context
		_ = config.LoadActiveConfigByName(cfg, cfg.Config.CurrentContext)
		f(*cfg)
	}
}

func HarnessAnchorfilesTestRepo(ctx common.Context) {
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath()
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	anchorfilesPathTest := fmt.Sprintf("%s/test/data/anchorfiles", repoRootPath)
	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(anchorfilesPathTest)
}
