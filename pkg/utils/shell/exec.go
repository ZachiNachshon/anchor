package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	gotty "github.com/mattn/go-tty"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ShellType string

const (
	bash ShellType = "/bin/bash"
	sh   ShellType = "/bin/sh"
)

type shellExecutor struct {
	shellType ShellType
}

func New() Shell {
	return &shellExecutor{
		shellType: sh, // is bin/sh enough?
	}
}

func (s *shellExecutor) ExecuteScriptRealtimeWithOutput(dir string, relativeScriptPath string, args ...string) (string, error) {
	path := fmt.Sprintf("%s/%s", dir, relativeScriptPath)
	// Args must include the command as Args[0]
	slice := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   slice,
		Dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	// Execute
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Collect std out script execution logs
	builder := strings.Builder{}
	oneByte := make([]byte, 100)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			logger.Warning("failed to read script output from stdout")
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		builder.WriteString(string(line))
	}

	cmd.Wait()
	return builder.String(), nil
}

func (s *shellExecutor) ExecuteScript(dir string, relativeScriptPath string, args ...string) error {
	path := fmt.Sprintf("%s/%s", dir, relativeScriptPath)
	// Args must include the command as Args[0]
	slice := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   slice,
		Dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	// Execute
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (s *shellExecutor) ExecuteTTY(script string) error {
	tty, err := gotty.Open()
	if err != nil {
		return err
	}
	defer tty.Close()

	args := strings.Fields(script)
	cmd := exec.Command(args[0], args[1:]...)

	// Setup the command's standard input/output/error
	cmd.Stdin = tty.Input()
	cmd.Stdout = tty.Output()
	cmd.Stderr = tty.Output()

	// Execute
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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

func (s *shellExecutor) ExecuteSilently(script string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *shellExecutor) ExecuteInBackground(script string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)
	// Temporary prevent logs verbosity from background process
	//cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Infof("Starting background process, PID: %d", cmd.Process.Pid)
	return nil
}
