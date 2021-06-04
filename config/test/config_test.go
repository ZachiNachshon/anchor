package test

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/kits"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestAuthor = "author: Dummy Name <dummy.name@gmail.com>"
var TestLicense = "license: TestsLicense"

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
			Name: "resolve local test anchorfiles path successfully",
			Func: ResolveLocalTestAnchorfilesPathSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ResolveLocalTestAnchorfilesPathSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := kits.GetDefaultTestConfigText()
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
			var items = kits.TemplateItems{
				Author:  TestAuthor,
				License: TestLicense,
			}
			yamlConfigText := kits.GetCustomTestConfigText(items)
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				nonViperConfig := converters.YamlToConfigObj(yamlConfigText)
				assert.EqualValues(t, cfg.Author, nonViperConfig.Author)
				assert.EqualValues(t, cfg.License, nonViperConfig.License)
				assert.EqualValues(t, cfg.Config, nonViperConfig.Config)
			})
		})
	})
}

var ResolveConfigWithDefaultsFromYamlTextSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			// Omit config items that should get default values
			yamlConfigText := kits.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				nonViperConfig := converters.YamlToConfigObj(yamlConfigText)
				assert.EqualValues(t, cfg.Author, config.DefaultAuthor)
				assert.EqualValues(t, cfg.License, config.DefaultLicense)
				assert.EqualValues(t, cfg.Config, nonViperConfig.Config)
			})
		})
	})
}
