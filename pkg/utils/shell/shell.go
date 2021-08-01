package shell

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/creack/pty"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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

func (s *shellExecutor) ExecuteScriptFile(dir string, relativeScriptPath string, args ...string) error {
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

func (s *shellExecutor) ExecuteScriptFileWithOutputToFile(
	workingDirectory string,
	relativeScriptPath string,
	outputFilePath string,
	args ...string) error {

	path := fmt.Sprintf("%s/%s", workingDirectory, relativeScriptPath)
	// Args must include the command as Args[0]
	slice := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path: path,
		Args: slice,
		Dir:  workingDirectory,
	}

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var _, stderrBuf bytes.Buffer
	// Script execution sends output to stderr instead of stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, file)

	err = cmd.Run()
	if err != nil {
		errStr := string(stderrBuf.Bytes())
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), errStr)
	}
	return nil
}

func (s *shellExecutor) Execute(script string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *shellExecutor) ExecuteWithOutputToFile(script string, outputFilePath string) error {
	cmd := exec.Command(string(s.shellType), "-c", script)

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var _, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, file)
	cmd.Stderr = io.MultiWriter(os.Stderr, file)

	err = cmd.Run()
	if err != nil {
		errStr := string(stderrBuf.Bytes())
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), errStr)
	}
	return nil
}

// ExecuteTTY example was inspired by - https://github.com/creack/pty#shell
func (s *shellExecutor) ExecuteTTY(script string) error {
	c := exec.Command(string(s.shellType), "-c", script)

	// Start the command with a pty
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end
	defer func() { _ = ptmx.Close() }() // Best effort

	// Handle pty size
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				logger.Debugf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done

	// Set stdin in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort

	// Copy stdin to the pty and the pty to stdout
	// NOTE: The goroutine will keep reading until the next keystroke before returning
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func (s *shellExecutor) ExecuteReturnOutput(script string) (string, error) {
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

func (s *shellExecutor) ClearScreen() error {
	// TODO: if windows support should be available in the future,
	//       adjust to cls by verifying os first
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
