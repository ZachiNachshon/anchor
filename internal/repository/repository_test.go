package repository

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/repository/local"
	"github.com/ZachiNachshon/anchor/internal/repository/remote"
	"github.com/ZachiNachshon/anchor/pkg/printer"

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
		{
			Name: "fail to resolve local or remote repositories origin ",
			Func: FailToResolveLocalOrRemoteRepositoriesOrigin,
		},
	}
	harness.RunTests(t, tests)
}

var FailToGetResolverDueToInvalidConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			res, err := GetRepositoryOriginByConfig(ctx, nil)
			assert.Nil(t, res, "expected invalid resolver")
			assert.NotNil(t, err, "expected to fail getting a repository resolver")
			assert.Contains(t, err.Error(), "missing required config value")
		})
	})
}

var GetLocalResolverFromConfigSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: https://github.com/ZachiNachshon/dummy-repo.git      
            revision: l33tf4k3c0mm1757r1n6 
            branch: some-branch
            clonePath: /best/path/ever
          local:
            path: /local/path/wins
`
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				res, err := GetRepositoryOriginByConfig(ctx, cfg.Config.ActiveContext.Context.Repository)
				assert.Equal(t, fmt.Sprintf("%T", local.NewLocalRepository(nil)), fmt.Sprintf("%T", res), "expected a local resolver")
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
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: https://github.com/ZachiNachshon/dummy-repo.git      
            revision: l33tf4k3c0mm1757r1n6 
            branch: some-branch
            clonePath: /best/path/ever
          local:
            path: ""
`
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				ctx.Registry().Set(printer.Identifier, printer.CreateFakePrinter())
				res, err := GetRepositoryOriginByConfig(ctx, cfg.Config.ActiveContext.Context.Repository)
				assert.Equal(t, fmt.Sprintf("%T", remote.NewRemoteRepository(nil)), fmt.Sprintf("%T", res), "expected a remote resolver")
				assert.Nil(t, err, "expected getting a resolver successfully")
			})
		})
	})
}

var FailToResolveLocalOrRemoteRepositoriesOrigin = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
   - name: test-cfg-ctx
     context:
       repository:
         remote:
         local:
           path: ""
`
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				res, err := GetRepositoryOriginByConfig(ctx, cfg.Config.ActiveContext.Context.Repository)
				assert.NotNil(t, err, "expected to fail")
				assert.Nil(t, res, "expected not to have a response object")
				assert.Contains(t, err.Error(), "could not resolve anchorfiles local repository path or git tracked remote repository")
			})
		})
	})
}
