package shell

import (
	"bytes"
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
	var output string
	if out, err := exec.Command(string(s.shellType), "-c", script).Output(); err != nil {
		return "", err
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
		//logger.Fatalf("Failure occurred: %v", err.Error())
	}

	//log.Println(stdBuffer.String())
	return nil
}
