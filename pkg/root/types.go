package root

import (
	"github.com/ZachiNachshon/anchor/common"
)

type RootCommandActions struct {
	LoadRepoOrFail func(ctx common.Context)
}

func DefineRootCommandActions() *RootCommandActions {
	return &RootCommandActions{
		LoadRepoOrFail: StartRootCommandLoadRepoOrFailFlow,
	}
}
