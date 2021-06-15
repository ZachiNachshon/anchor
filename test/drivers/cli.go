package drivers

import (
	"bytes"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"io/ioutil"
)

type CLIActions interface {
	RunCommand(cmd models.AnchorCommand, args ...string) (string, error)
}

type cliRunnerImpl struct {
	CLIActions
}

func CLI() CLIActions {
	return &cliRunnerImpl{}
}

func (cli *cliRunnerImpl) RunCommand(cmd models.AnchorCommand, args ...string) (string, error) {
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
