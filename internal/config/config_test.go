package config

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var TestAuthor = "author: Dummy Name <dummy.name@gmail.com>"
var TestLicense = "license: TestsLicense"
var TestClonePath = "/test/clone/path"

func Test_ConfigShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "resolve config from YAML text successfully",
			Func: ResolveConfigFromYamlTextSuccessfully,
		},
		{
			Name: "resolve config with defaults from YAML text successfully",
			Func: ResolveConfigWithDefaultsFromYamlTextSuccessfully,
		},
		{
			Name: "resolve local anchorfiles test repo successfully",
			Func: ResolveLocalAnchorfilesTestRepoSuccessfully,
		},
		//{
		//	Name: "fail to resolve local repository due to missing config",
		//	Func: FailToResolveLocalRepositoryDueToMissingConfig,
		//},
		//{
		//	Name: "fail to resolve remote repository due to invalid config",
		//	Func: FailToResolveRemoteRepositoryDueToInvalidConfig,
		//},
	}
	harness.RunTests(t, tests)
}

var ResolveLocalAnchorfilesTestRepoSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			yamlConfigText := GetDefaultTestConfigText()
			withConfig(ctx, yamlConfigText, func(config *AnchorConfig) {
				harnessAnchorfilesTestRepo(&ctx)
				assert.DirExists(t, ctx.AnchorFilesPath(),
					"cannot resolve anchorfiles test repo. path: %s", ctx.AnchorFilesPath())
			})
		})
	})
}

var ResolveConfigFromYamlTextSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			// Override default values explicitly
			var items = TemplateItems{
				Author:                        TestAuthor,
				License:                       TestLicense,
				CurrentContext:                "1st-anchorfiles",
				FirstContextName:              "1st-anchorfiles",
				FirstContextClonePath:         TestClonePath,
				FirstContextRemoteRepoBranch:  "1st-test-branch",
				SecondContextName:             "2nd-anchorfiles",
				SecondContextClonePath:        TestClonePath,
				SecondContextRemoteRepoBranch: "2nd-test-branch",
			}
			yamlConfigText := GetCustomTestConfigText(items)
			withConfig(ctx, yamlConfigText, func(cfg *AnchorConfig) {
				cfg.Config.ActiveContext = nil
				nonViperConfig := YamlToConfigObj(yamlConfigText)
				assert.EqualValues(t, nonViperConfig.Author, cfg.Author)
				assert.EqualValues(t, nonViperConfig.License, cfg.License)
				assert.EqualValues(t, nonViperConfig.Config, cfg.Config)
			})
		})
	})
}

var ResolveConfigWithDefaultsFromYamlTextSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			// Omit config items that should get default values
			yamlConfigText := GetDefaultTestConfigText()
			withConfig(ctx, yamlConfigText, func(cfg *AnchorConfig) {
				cfgCtxName := "1st-anchorfiles"
				cfgManager := NewManager()
				defaultClonePath, _ := cfgManager.GetDefaultRepoClonePath(cfgCtxName)
				nonViperConfig := YamlToConfigObj(yamlConfigText)
				assert.NotNil(t, nonViperConfig, "expected a valid config object")
				assert.EqualValues(t, DefaultAuthor, cfg.Author)
				assert.EqualValues(t, DefaultLicense, cfg.License)
				assert.EqualValues(t, defaultClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath)
			})
		})
	})
}

//var FailToResolveRemoteRepositoryDueToInvalidConfig = func(t *testing.T) {
//	with.Context(func(ctx common.Context) {
//		with.Logging(ctx, t, func(logger logger.Logger) {
//			yamlConfigText := `
//config:
// currentContext: test-cfg-ctx
// contexts:
//   - name: test-cfg-ctx
//     context:
//       repository:
//         remote:
//`
//			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
//				verifyConfigCallCount := 0
//				fakeRemoteActions := CreateFakeRemoteActions()
//				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
//					verifyConfigCallCount++
//					return fmt.Errorf("invalid config")
//				}
//
//				rslvr := &RemoteResolver{
//					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
//					RemoteActions: fakeRemoteActions,
//				}
//				repoPath, err := rslvr.Load(ctx)
//				assert.NotNil(t, err, "expected to fail on remote resolver")
//				assert.Equal(t, "invalid config", err.Error())
//				assert.Equal(t, 1, verifyConfigCallCount)
//				assert.Equal(t, "", repoPath, "expected not to have a repository path")
//			})
//		})
//	})
//}
//
//var FailToResolveLocalRepositoryDueToMissingConfig = func(t *testing.T) {
//	with.Context(func(ctx common.Context) {
//		with.Logging(ctx, t, func(logger logger.Logger) {
//			yamlConfigText := `
//config:
// currentContext: test-cfg-ctx
// contexts:
//   - name: test-cfg-ctx
//     context:
//       repository:
//         local:
//`
//			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
//				rslvr := &LocalResolver{
//					LocalConfig: cfg.Config.ActiveContext.Context.Repository.Local,
//				}
//				repoPath, err := rslvr.Load(ctx)
//				assert.NotNil(t, err, "expected to fail on local resolver")
//				assert.Equal(t, "invalid local repository configuration", err.Error())
//				assert.Equal(t, "", repoPath, "expected not to have a repository path")
//			})
//		})
//	})
//}

// Duplicate these methods from the with.go test utility to prevent 'import cycle not allowed in test'
// config package is the only exception for this use case
func withContext(f func(ctx common.Context)) {
	// Every context must have a new registry instance
	ctx := common.EmptyAnchorContext(registry.New())
	f(ctx)
}

func withConfig(ctx common.Context, content string, f func(config *AnchorConfig)) {
	cfgManager := NewManager()
	if err := cfgManager.SetupConfigInMemoryLoader(content); err != nil {
		logger.Fatalf("Failed to create a fake config loader. error: %s", err)
	} else {
		cfg, err := cfgManager.CreateConfigObject()
		if err != nil {
			logger.Fatalf("Failed to create a fake config loader. error: %s", err)
		}
		SetInContext(ctx, cfg)
		// set current config context as the active config context
		_ = cfgManager.SwitchActiveConfigContextByName(cfg, cfg.Config.CurrentContext)
		ctx.Registry().Set(Identifier, cfgManager)
		f(cfg)
	}
}

func withLogging(ctx common.Context, t *testing.T, verbose bool, f func(logger logger.Logger)) {
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

func harnessAnchorfilesTestRepo(ctx *common.Context) {
	path, _ := ioutils.GetWorkingDirectory()
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath(path)
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	anchorfilesPathTest := fmt.Sprintf("%s/test/data/anchorfiles", repoRootPath)
	(*ctx).(common.AnchorFilesPathSetter).SetAnchorFilesPath(anchorfilesPathTest)
}
