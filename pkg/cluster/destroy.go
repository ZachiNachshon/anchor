package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
)

func Destroy() error {
	name := common.GlobalOptions.KindClusterName
	logger.PrintWarning(fmt.Sprintf("Remember to backup all your data!"))
	in := input.NewYesNoInput()
	destroyInputFormat := fmt.Sprintf("Are you sure you want to destroy Kubernetes cluster [%v]?", name)

	if result, err := in.WaitForInput(destroyInputFormat); err != nil || !result {
		logger.Info("skipping.")
	} else {

		logger.PrintCommandHeader(fmt.Sprintf("Destroying cluster %v", name))

		// Kill possible running kubectl proxy
		_ = KillAllRunningKubectl()

		// Remove docker registry and clean /etc/host entry
		_ = DeleteRegistry()

		destroyCmd := fmt.Sprintf("kind delete cluster --name %v", name)
		if err := common.ShellExec.Execute(destroyCmd); err != nil {
			return err
		}
	}
	return nil
}
