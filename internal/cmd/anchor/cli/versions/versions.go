package versions

import (
	"github.com/ZachiNachshon/anchor/internal/common"
)

type CliVersionsFunc func(ctx common.Context) error

var CliVersions = func(ctx common.Context) error {
	return nil
}
