package shell

import (
	"bufio"
	"os/exec"

	"github.com/kit/pkg/logger"
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

func (s *shellExecutor) ExecShellWithOutput(command string) (string, error) {
	var output string
	if out, err := exec.Command(string(s.shellType), "-c", command).Output(); err != nil {
		return "", err
	} else {
		output = string(out)
	}
	return output, nil
}

func (s *shellExecutor) ExecShell(command string) {
	cmd := exec.Command(string(s.shellType), "-c", command)

	logger.Info("\n")

	stdout, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	oneByte := make([]byte, 100)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			// Do nothing, no need to print EOF
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		logger.Info(string(line))
	}

	err := cmd.Wait()
	if err != nil {
		logger.Fatalf("Failure occurred: %v", err.Error())
	}
}
