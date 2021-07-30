package test

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
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
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				harness.HarnessAnchorfilesTestRepo(ctx)
				assert.True(t, ioutils.IsValidPath(ctx.AnchorFilesPath()),
					"cannot resolve anchorfiles test repo. path: %s", ctx.AnchorFilesPath())
			})
		})
	})
}

var ResolveConfigFromYamlTextSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			// Override default values explicitly
			var items = config.TemplateItems{
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
			yamlConfigText := config.GetCustomTestConfigText(items)
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.ActiveContext = nil
				nonViperConfig := converters.YamlToConfigObj(yamlConfigText)
				assert.EqualValues(t, nonViperConfig.Author, cfg.Author)
				assert.EqualValues(t, nonViperConfig.License, cfg.License)
				assert.EqualValues(t, nonViperConfig.Config, cfg.Config)
			})
		})
	})
}

var ResolveConfigWithDefaultsFromYamlTextSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			// Omit config items that should get default values
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfgCtxName := "1st-anchorfiles"
				defaultClonePath, _ := config.GetDefaultRepoClonePath(cfgCtxName)
				nonViperConfig := converters.YamlToConfigObj(yamlConfigText)
				assert.NotNil(t, nonViperConfig, "expected a valid config object")
				assert.EqualValues(t, config.DefaultAuthor, cfg.Author)
				assert.EqualValues(t, config.DefaultLicense, cfg.License)
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
//  currentContext: test-cfg-ctx
//  contexts:
//    - name: test-cfg-ctx
//      context:
//        repository:
//          remote:
//`
//			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
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
//				repoPath, err := rslvr.ResolveRepository(ctx)
//				assert.NotNil(t, err, "expected to fail on remote resolver")
//				assert.Equal(t, "invalid config", err.Error())
//				assert.Equal(t, 1, verifyConfigCallCount)
//				assert.Equal(t, "", repoPath, "expected not to have a repository path")
//			})
//		})
//	})
//}

//var FailToResolveLocalRepositoryDueToMissingConfig = func(t *testing.T) {
//	with.Context(func(ctx common.Context) {
//		with.Logging(ctx, t, func(logger logger.Logger) {
//			yamlConfigText := `
//config:
//  currentContext: test-cfg-ctx
//  contexts:
//    - name: test-cfg-ctx
//      context:
//        repository:
//          local:
//`
//			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
//				rslvr := &LocalResolver{
//					LocalConfig: cfg.Config.ActiveContext.Context.Repository.Local,
//				}
//				repoPath, err := rslvr.ResolveRepository(ctx)
//				assert.NotNil(t, err, "expected to fail on local resolver")
//				assert.Equal(t, "invalid local repository configuration", err.Error())
//				assert.Equal(t, "", repoPath, "expected not to have a repository path")
//			})
//		})
//	})
//}
