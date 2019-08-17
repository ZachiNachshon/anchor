package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func LogRunningPod(identifier string, namespace string) error {
	var name string
	var err error

	if name, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	podsParser := NewPodsSelector()
	if err := podsParser.PrepareOptions(name, namespace); err != nil {
		return err
	}

	if len(podsParser.podsInfo) == 0 {
		msg := fmt.Sprintf("\n  No pod(s) could be found for name: %v, namespace: %v.", name, namespace)
		logger.Info(msg)
		return nil
	}

	if podInfo, err := podsParser.SelectPod(name, namespace); err != nil {
		return err
	} else {
		logPodCmd := fmt.Sprintf("kubectl logs -f %v -n %v", podInfo.Name, common.GlobalOptions.KindClusterName)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + logPodCmd + "\n")
		}

		logger.PrintCommandHeader(fmt.Sprintf("Logging %v", name))
		if err := common.ShellExec.ExecuteTTY(logPodCmd); err != nil {
			return err
		}
	}

	return nil
}
