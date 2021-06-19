package controller

import "github.com/ZachiNachshon/anchor/common"

type ControllerActions struct {
	Install   func(ctx common.Context) error
	Uninstall func(ctx common.Context) error
	List      func(ctx common.Context) error
	Status    func(ctx common.Context) error
	Versions  func(ctx common.Context) error
}

func DefineControllerActions() *ControllerActions {
	return &ControllerActions{}
}
