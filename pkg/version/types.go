package version

import "github.com/ZachiNachshon/anchor/common"

type VersionActions struct {
	Version func(ctx common.Context) error
}

func DefineVersionActions() *VersionActions {
	return &VersionActions{
		Version: StartVersionVersionFlow,
	}
}
