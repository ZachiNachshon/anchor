package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

func (lr *LocalResolver) ResolveRepository(ctx common.Context) (string, error) {
	if lr.LocalConfig != nil {
		pathToUse := lr.LocalConfig.Path
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
