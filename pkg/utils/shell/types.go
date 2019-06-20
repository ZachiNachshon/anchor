package shell

type Shell interface {
	ExecShellWithOutput(text string) (string, error)
	ExecShell(text string)
}
