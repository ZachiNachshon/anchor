package zsh

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ZshCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "generate zsh completion successfully",
			Func: GenerateZshCompletionSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var GenerateZshCompletionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				var rootCmd = &cobra.Command{
					Use:   "anchor",
					Short: "root cmd",
					Long:  `root cmd`,
				}
				_, err := drivers.CLI().RunCommand(NewCommand(rootCmd))
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}
