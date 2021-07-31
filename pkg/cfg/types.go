package cfg

import "github.com/ZachiNachshon/anchor/common"

type ConfigurationActions struct {
	Print           func(ctx common.Context) error
	Edit            func(ctx common.Context) error
	UseContext      func(ctx common.Context, cfgCtxName string) error
	SetContextEntry func(ctx common.Context, cfgCtxName string, changes map[string]string) error
}

func DefineConfigurationActions() *ConfigurationActions {
	return &ConfigurationActions{
		Print:           StartConfigPrintFlow,
		Edit:            StartConfigEditFlow,
		UseContext:      StartConfigUseContextFlow,
		SetContextEntry: StartConfigSetContextEntryFlow,
	}
}
