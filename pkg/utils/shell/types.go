package shell

type Shell interface {
	ExecuteScript(dir string, relativeScriptPath string, args ...string) error
	ExecuteWithOutput(script string) (string, error)
	Execute(script string) error
	ExecuteTTY(script string) error
	ExecuteInBackground(script string) error
}
