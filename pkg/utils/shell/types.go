package shell

type Shell interface {
	ExecShellWithOutput(script string) (string, error)
	ExecShell(script string) error
}

type Installer interface {
	install() error
	verify() error
	Check() error
}
