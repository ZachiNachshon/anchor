package install

import (
	"github.com/ZachiNachshon/anchor/internal/common"
)

type ControllerInstallFunc func(ctx common.Context) error

var ControllerInstall = func(ctx common.Context) error {
	return nil
}
