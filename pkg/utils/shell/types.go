package shell

type Shell interface {
	ExecuteWithOutput(script string) (string, error)
	Execute(script string) error
	ExecuteTTY(script string) error
	ExecuteInBackground(script string) error
}
