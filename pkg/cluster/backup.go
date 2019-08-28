package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func Backup(identifier string, namespace string) (string, error) {
	var name string
	var err error

	if name, err = locator.DirLocator.Name(identifier); err != nil {
		return "", err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Backing up %v", name))
	}

	if err := unMountHostPath(name, namespace, true); err != nil {
		return "", err
	}
	return "", nil
}
