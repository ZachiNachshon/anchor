package drivers

import (
	"bytes"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/spf13/cobra"
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
	//cmd.InitFlags()
	//cmd.InitSubCommands()

	cobraCmd := cmd.GetCobraCmd()
	b := bytes.NewBufferString("")
	cobraCmd.SetOut(b)
	if args != nil {
		cobraCmd.SetArgs(args)
	} else {
		cobraCmd.SetArgs([]string{})
	}
	// In order to catch all run phases we must have Run/RunE implemented on the cmd
	//  * PersistentPreRun()
	//  * PreRun()
	//  * Run()
	//  * PostRun()
	//  * PersistentPostRun()
	if cobraCmd.Run == nil {
		cobraCmd.Run = func(cmd *cobra.Command, args []string) {}
	}

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
