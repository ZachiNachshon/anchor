package common

var GlobalOptions = CmdRootOptions{
	Verbose: false,
}

type CmdRootOptions struct {
	ConfigFile string

	// Log options
	Verbose bool

	// Additional Params
}
