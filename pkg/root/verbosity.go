package root

import (
	"github.com/ZachiNachshon/anchor/logger"
)

func StartRootCommandVerbositySetterFlow(l logger.Logger, verbose bool) error {
	level := "info"
	if verbose {
		level = "debug"
	}
	if err := l.SetVerbosityLevel(level); err != nil {
		return err
	}
	return nil
}
