package app

import "github.com/ZachiNachshon/anchor/common"

type ApplicationActions struct {
	Select func(ctx common.Context) error
	Status func(ctx common.Context) error
}

func DefineApplicationActions() *ApplicationActions {
	return &ApplicationActions{
		Select: StartApplicationSelectionFlow,
		Status: StartApplicationStatusFlow,
	}
}
