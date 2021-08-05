package controller

import "github.com/ZachiNachshon/anchor/common"

type ControllerActions struct {
	Install func(ctx common.Context) error
}

func DefineControllerActions() *ControllerActions {
	return &ControllerActions{
		Install: StartControllerInstallFlow,
	}
}
