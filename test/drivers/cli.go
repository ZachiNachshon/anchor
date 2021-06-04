package drivers

import (
	"bytes"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// TODO: move to prod code
type commandInTest interface {
	GetCobraCmd() *cobra.Command
	InitFlags()
	InitSubCommands()
}

var (
	CLI = &cliRunnerImpl{}
)

type cliRunnerImpl struct{}

func (cli *cliRunnerImpl) RunCommand(cmd commandInTest, args ...string) (string, error) {
	cmd.InitFlags()
	cmd.InitSubCommands()

	cobraCmd := cmd.GetCobraCmd()
	b := bytes.NewBufferString("")
	cobraCmd.SetOut(b)
	cobraCmd.SetArgs(args)
	err := cobraCmd.Execute()
	if err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return string(out), nil
}
