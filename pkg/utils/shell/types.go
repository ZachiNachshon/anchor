package shell

type Shell interface {
	ExecuteWithOutput(script string) (string, error)
	Execute(script string) error
	ExecuteInBackground(script string) error
}

type Installer interface {
	install() error
	verify() error
	Check() error
}
