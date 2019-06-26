package shell

import (
	"bytes"
	"github.com/anchor/pkg/logger"
	"io"
	"os"
	"os/exec"
)

type ShellType string

const (
	BASH ShellType = "/bin/bash"
	SH   ShellType = "/bin/sh"
)

type shellExecutor struct {
	shellType ShellType
}

func NewShellExecutor(sType ShellType) Shell {
	return &shellExecutor{
		shellType: sType,
	}
}

func (s *shellExecutor) ExecuteWithOutput(script string) (string, error) {
	cmd := exec.Command(string(s.shellType), "-c", script)

	var output string
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if out, err := cmd.Output(); err != nil {
		return stderr.String(), err
	} else {
		output = string(out)
	}
	return output, nil
}

func (s *shellExecutor) Execute(script string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *shellExecutor) ExecuteInBackground(script string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		// TODO: Change to warn once implemented
		logger.Info(err.Error())
		return err
	}
	logger.Infof("Starting background process, PID: %d", cmd.Process.Pid)
	return nil
}
