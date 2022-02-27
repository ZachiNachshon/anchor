package shell

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"path/filepath"
	"strings"

	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/creack/pty"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const (
	Identifier string = "shell"
)

type Shell interface {
	ExecuteScriptFile(relativeScriptPath string, args ...string) error
	ExecuteScriptFileWithOutputToFile(
		relativeScriptPath string,
		outputFilePath string,
		args ...string) error

	ExecuteScriptFileSilentlyWithOutputToFile(
		relativeScriptPath string,
		outputFilePath string,
		args ...string) error

	Execute(script string) error
	ExecuteWithOutputToFile(script string, outputFilePath string) error
	ExecuteSilentlyWithOutputToFile(script string, outputFilePath string) error

	ExecuteReturnOutput(script string) (string, error)
	ExecuteSilently(script string) error

	// ExecuteTTY executes a script as a TeleTYpewrite, this allow us to run interactive commands
	// by creating a unix pseudo-terminals.
	// Used for cases of starting applications programmatically such as vi, sublime etc..
	ExecuteTTY(script string) error
	ExecuteInBackground(script string) error
	ClearScreen() error
}

type ShellType string

const (
	bash ShellType = "/bin/bash"
	sh   ShellType = "/bin/sh"
)

type shellExecutor struct {
	shellType ShellType
	ctx       common.Context
}

func New(ctx common.Context) Shell {
	return &shellExecutor{
		ctx:       ctx,
		shellType: sh, // is bin/sh enough?
	}
}

func (s *shellExecutor) ExecuteScriptFile(relativeScriptPath string, args ...string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	path := fmt.Sprintf("%s/%s", workingDirectory, relativeScriptPath)
	// Args must be appended with the command as Args[0]
	slice := append([]string{path}, args...)

	envs := getEnvVarsForScriptFileExec(path, workingDirectory)
	cmd := &exec.Cmd{
		Path:   path,
		Args:   slice,
		Env:    envs,
		Dir:    workingDirectory,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
		Stdin:  os.Stdin,
	}

	// Execute
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (s *shellExecutor) ExecuteScriptFileWithOutputToFile(
	relativeScriptPath string,
	outputFilePath string,
	args ...string) error {

	workingDirectory := s.ctx.AnchorFilesPath()
	path := fmt.Sprintf("%s/%s", workingDirectory, relativeScriptPath)
	// Args must be appended with the command as Args[0]
	slice := append([]string{path}, args...)

	envs := getEnvVarsForScriptFileExec(path, workingDirectory)
	cmd := &exec.Cmd{
		Path: path,
		Args: slice,
		Env:  envs,
		Dir:  workingDirectory,
	}

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf, file)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf, file)
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		errStr := extractLastErrorLine(string(stderrBuf.Bytes()))
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), strings.TrimSpace(errStr))
	}
	return nil
}

func (s *shellExecutor) ExecuteScriptFileSilentlyWithOutputToFile(
	relativeScriptPath string,
	outputFilePath string,
	args ...string) error {

	workingDirectory := s.ctx.AnchorFilesPath()
	path := fmt.Sprintf("%s/%s", workingDirectory, relativeScriptPath)
	// Args must be appended with the command as Args[0]
	slice := append([]string{path}, args...)

	envs := getEnvVarsForScriptFileExec(path, workingDirectory)
	cmd := &exec.Cmd{
		Path: path,
		Args: slice,
		Env:  envs,
		Dir:  workingDirectory,
	}

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdoutBuf, file)
	cmd.Stderr = io.MultiWriter(&stderrBuf, file)

	err = cmd.Run()
	if err != nil {
		errStr := extractLastErrorLine(string(stderrBuf.Bytes()))
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), strings.TrimSpace(errStr))
	}
	return nil
}

func (s *shellExecutor) Execute(script string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *shellExecutor) ExecuteWithOutputToFile(script string, outputFilePath string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var _, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stderrBuf, file)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf, file)
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		errStr := extractLastErrorLine(string(stderrBuf.Bytes()))
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), strings.TrimSpace(errStr))
	}
	return nil
}

func (s *shellExecutor) ExecuteSilentlyWithOutputToFile(script string, outputFilePath string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory

	file, err := ioutils.CreateOrOpenFile(outputFilePath)
	if err != nil {
		return err
	}

	var _, stderrBuf bytes.Buffer
	cmd.Stdout = file
	cmd.Stderr = io.MultiWriter(&stderrBuf, file)

	err = cmd.Run()
	if err != nil {
		errStr := extractLastErrorLine(string(stderrBuf.Bytes()))
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), strings.TrimSpace(errStr))
	}
	return nil
}

// ExecuteTTY example was inspired by - https://github.com/creack/pty#shell
func (s *shellExecutor) ExecuteTTY(script string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	c := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	c.Dir = workingDirectory

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
			if err = pty.InheritSize(os.Stdin, ptmx); err != nil {
				logger.Errorf("error resizing pty: %s", err)
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
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory

	var output string
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	if out, err := cmd.Output(); err != nil {
		return stderr.String(), err
	} else {
		output = string(out)
	}
	return output, nil
}

func (s *shellExecutor) ExecuteSilently(script string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory

	var _, stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error: %s, stderr: %s", err.Error(), strings.TrimSpace(stderrBuf.String()))
	}
	return nil
}

func (s *shellExecutor) ExecuteInBackground(script string) error {
	workingDirectory := s.ctx.AnchorFilesPath()
	substituteEnvVarScript := expandScriptEnvVars(script)
	cmd := exec.Command(string(s.shellType), "-c", substituteEnvVarScript)
	cmd.Dir = workingDirectory
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

func extractLastErrorLine(errStr string) string {
	split := strings.SplitAfterN(errStr, "\n", 2)
	if split != nil && len(split) > 1 {
		return split[1]
	}
	return errStr
}

func expandScriptEnvVars(script string) string {
	return os.ExpandEnv(script)
}

// Append required env vars for scripted languages in order to
// "lock" the working directory as the execution path
func getEnvVarsForScriptFileExec(scriptFile string, workingDirectory string) []string {
	envs := os.Environ()
	if isPythonScript(scriptFile) {
		// To allow Python to resolve imports properly, we're setting the
		// PYTHONPATH env var with the exec path which is the working dir
		pythonPath := fmt.Sprintf("PYTHONPATH=%s", workingDirectory)
		envs = append(envs, pythonPath)
	}
	return envs
}

func isPythonScript(scriptFile string) bool {
	return filepath.Ext(scriptFile) == ".py"
}
