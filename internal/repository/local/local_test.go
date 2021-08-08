package local

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RemoteShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "resolve local repository successfully",
			Func: ResolveLocalRepositorySuccessfully,
		},
		{
			Name: "fail to resolve local repository due to invalid path",
			Func: FailToResolveLocalRepositoryDueToInvalidPath,
		},
	}
	harness.RunTests(t, tests)
}

var ResolveLocalRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.HarnessAnchorfilesTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          local:
            path: %s
`, ctx.AnchorFilesPath())
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				repo := &LocalRepository{
					LocalConfig: cfg.Config.ActiveContext.Context.Repository.Local,
				}
				repoPath, err := repo.Load(ctx)
				assert.Nil(t, err, "expected to succeed on local resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected a valid local repository path")
			})
		})
	})
}

var FailToResolveLocalRepositoryDueToInvalidPath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          local:
            path: /invalid/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				repo := &LocalRepository{
					LocalConfig: cfg.Config.ActiveContext.Context.Repository.Local,
				}
				repoPath, err := repo.Load(ctx)
				assert.NotNil(t, err, "expected to fail on local resolver")
				assert.Contains(t, err.Error(), "local anchorfiles repository path is invalid")
				assert.Equal(t, "", repoPath, "expected not to have a repository path path")
			})
		})
	})
}
