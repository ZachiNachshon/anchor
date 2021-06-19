package app

import "github.com/ZachiNachshon/anchor/common"

type ApplicationActions struct {
	Install   func(ctx common.Context) error
	Uninstall func(ctx common.Context) error
	List      func(ctx common.Context) error
	Status    func(ctx common.Context) error
	Versions  func(ctx common.Context) error
}

func DefineApplicationActions() *ApplicationActions {
	return &ApplicationActions{
		Install: StartApplicationInstallFlow,
		List:    StartApplicationListFlow,
	}
}
