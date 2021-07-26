package app

import "github.com/ZachiNachshon/anchor/common"

type ApplicationActions struct {
	Select   func(ctx common.Context) error
	List     func(ctx common.Context) error
	Status   func(ctx common.Context) error
	Versions func(ctx common.Context) error
}

func DefineApplicationActions() *ApplicationActions {
	return &ApplicationActions{
		Select: StartApplicationSelectionFlow,
		List:   StartApplicationListFlow,
	}
}
