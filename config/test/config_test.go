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
var TestClonePath = "clonePath: /test/clone/path"

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
				Author:    TestAuthor,
				License:   TestLicense,
				ClonePath: TestClonePath,
			}
			yamlConfigText := config.GetCustomTestConfigText(items)
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
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
				defaultClonePath, _ := config.GetDefaultRepoClonePath()
				nonViperConfig := converters.YamlToConfigObj(yamlConfigText)
				assert.NotNil(t, nonViperConfig, "expected a valid config object")
				assert.EqualValues(t, config.DefaultAuthor, cfg.Author)
				assert.EqualValues(t, config.DefaultLicense, cfg.License)
				assert.EqualValues(t, defaultClonePath, cfg.Config.Repository.Remote.ClonePath)
			})
		})
	})
}
