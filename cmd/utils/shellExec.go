package utils

import (
	"bufio"
	"github.com/kit/cmd/logger"
	"os/exec"
)

func ExecShellWithOutput(command string) (string, error) {
	var output string
	// TODO: Need to allow other shells e.g /bin/sh for alpine
	if out, err := exec.Command("/bin/bash", "-c", command).Output(); err != nil {
		return "", err
	} else {
		output = string(out)
	}
	return output, nil
}

func ExecShell(command string) {
	// TODO: Need to allow other shells e.g /bin/sh for alpine
	cmd := exec.Command("/bin/bash", "-c", command)

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

	_ = cmd.Wait()
}
