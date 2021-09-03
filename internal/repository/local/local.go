package local

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

type LocalRepository interface {
	Load(ctx common.Context) (string, error)
}
type localRepositoryImpl struct {
	localConfig *config.Local
}

func NewLocalRepository(localConfig *config.Local) *localRepositoryImpl {
	return &localRepositoryImpl{
		localConfig: localConfig,
	}
}

func (lr *localRepositoryImpl) Load(ctx common.Context) (string, error) {
	if lr.localConfig != nil {
		pathToUse := lr.localConfig.Path
		if len(pathToUse) > 0 {
			if !ioutils.IsValidPath(pathToUse) {
				return "", fmt.Errorf("local anchorfiles repository path is invalid. path: %s", pathToUse)
			} else {
				return pathToUse, nil
			}
		}
	}

	return "", fmt.Errorf("invalid local repository configuration")
}
