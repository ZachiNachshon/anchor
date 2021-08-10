package bash

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BashCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "generate bash completion successfully",
			Func: GenerateBashCompletionSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var GenerateBashCompletionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				var rootCmd = &cobra.Command{
					Use:   "anchor",
					Short: "root cmd",
					Long:  `root cmd`,
				}
				command, err := NewCommand(rootCmd)
				_, err = drivers.CLI().RunCommand(command)
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}
