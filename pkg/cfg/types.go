package cfg

import "github.com/ZachiNachshon/anchor/common"

type ConfigurationActions struct {
	Print func(ctx common.Context) error
	Edit  func(ctx common.Context) error
}

func DefineConfigurationActions() *ConfigurationActions {
	return &ConfigurationActions{
		Print: StartConfigPrintFlow,
		Edit:  StartConfigEditFlow,
	}
}
