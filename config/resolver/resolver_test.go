package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ResolverShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail to get resolver due to invalid config",
			Func: FailToGetResolverDueToInvalidConfig,
		},
		{
			Name: "get local resolver from config successfully",
			Func: GetLocalResolverFromConfigSuccessfully,
		},
		{
			Name: "get remote resolver from config successfully",
			Func: GetRemoteResolverFromConfigSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailToGetResolverDueToInvalidConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Nil(t, res, "expected invalid resolver")
				assert.NotNil(t, err, "expected to fail getting a repository resolver")
				assert.Contains(t, err.Error(), "missing required config value")
			})
		})
	})
}

var GetLocalResolverFromConfigSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git      
      revision: l33tf4k3c0mm1757r1n6 
      branch: some-branch
      clonePath: /best/path/ever
    local:
      path: /local/path/wins
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Equal(t, fmt.Sprintf("%T", &LocalResolver{}), fmt.Sprintf("%T", res), "expected a local resolver")
				assert.Nil(t, err, "expected getting a resolver successfully")
			})
		})
	})
}

var GetRemoteResolverFromConfigSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git      
      revision: l33tf4k3c0mm1757r1n6 
      branch: some-branch
      clonePath: /best/path/ever
    local:
      path: ""
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Equal(t, fmt.Sprintf("%T", &RemoteResolver{}), fmt.Sprintf("%T", res), "expected a remote resolver")
				assert.Nil(t, err, "expected getting a resolver successfully")
			})
		})
	})
}
