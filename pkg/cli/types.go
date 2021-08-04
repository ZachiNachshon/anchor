package cli

import "github.com/ZachiNachshon/anchor/common"

type CliActions struct {
	Versions func(ctx common.Context) error
}

func DefineCliActions() *CliActions {
	return &CliActions{
		Versions: StartCliVersionsFlow,
	}
}
