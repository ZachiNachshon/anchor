package root

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
)

type RootCommandActions struct {
	LoadRepoOrFail     func(ctx common.Context)
	SetLoggerVerbosity func(l logger.Logger, verbose bool) error
}

func DefineRootCommandActions() *RootCommandActions {
	return &RootCommandActions{
		LoadRepoOrFail:     StartRootCommandLoadRepoOrFailFlow,
		SetLoggerVerbosity: StartRootCommandVerbositySetterFlow,
	}
}
