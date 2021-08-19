package with

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"os"
	"testing"
)

func Context(f func(ctx common.Context)) {
	// Every context must have a new registry instance
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
	if fakeLogger, err := logger.CreateFakeTestingLogger(t, verbose); err != nil {
		println("Failed to create a fake testing logger. error: %s", err)
		os.Exit(1)
	} else {
		fakeLoggerManager := logger.CreateFakeLoggerManager()
		fakeLoggerManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			return fakeLogger, nil
		}
		fakeLoggerManager.AppendFileLoggerMock = func(level string) (logger.Logger, error) {
			return fakeLogger, nil
		}
		fakeLoggerManager.SetVerbosityLevelMock = func(level string) error {
			resolveAdapter := fakeLogger.(logger.LoggerLogrusAdapter)
			err = resolveAdapter.SetStdoutVerbosityLevel(level)
			if err != nil {
				return err
			}

			err = resolveAdapter.SetFileVerbosityLevel(level)
			if err != nil {
				return err
			}
			return nil
		}

		fakeLoggerManager.GetDefaultLoggerLogFilePathMock = func() (string, error) {
			return "/testing/logger/path", nil
		}
		fakeLoggerManager.SetActiveLoggerMock = func(log *logger.Logger) error {
			return nil
		}

		err = fakeLoggerManager.SetActiveLogger(&fakeLogger)
		if err != nil {
			return
		}
		ctx.Registry().Set(logger.Identifier, fakeLoggerManager)
		f(fakeLogger)
	}
}

func Config(ctx common.Context, content string, f func(cfg *config.AnchorConfig)) {
	cfgManager := config.NewManager()
	if err := cfgManager.SetupConfigInMemoryLoader(content); err != nil {
		logger.Fatalf("Failed to create a fake config loader. error: %s", err)
	} else {
		cfg, err := cfgManager.CreateConfigObject()
		if err != nil {
			logger.Fatalf("Failed to create a fake config loader. error: %s", err)
		}
		config.SetInContext(ctx, cfg)
		// set current config context as the active config context
		_ = cfgManager.SwitchActiveConfigContextByName(cfg, cfg.Config.CurrentContext)
		ctx.Registry().Set(config.Identifier, cfgManager)
		f(cfg)
	}
}

func HarnessAnchorfilesTestRepo(ctx *common.Context) {
	path, _ := ioutils.GetWorkingDirectory()
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath(path)
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	anchorfilesPathTest := fmt.Sprintf("%s/test/data/anchorfiles", repoRootPath)
	(*ctx).(common.AnchorFilesPathSetter).SetAnchorFilesPath(anchorfilesPathTest)
}
