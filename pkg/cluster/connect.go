package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func ConnectToPod(identifier string, namespace string) error {
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
		execPodBashCmd := fmt.Sprintf("kubectl exec -it %v %v -n %v", podInfo.Name, shell.BASH, common.GlobalOptions.KindClusterName)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + execPodBashCmd + "\n")
		}

		logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v)", name, shell.BASH))
		if err := common.ShellExec.ExecuteTTY(execPodBashCmd); err != nil {

			// Fallback to /bin/sh if /bin/bash is not available
			execPodShCmd := fmt.Sprintf("kubectl exec -it %v %v -n %v", podInfo.Name, shell.SH, common.GlobalOptions.KindClusterName)
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + execPodShCmd + "\n")
			}

			logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v)", name, shell.SH))
			if err := common.ShellExec.ExecuteTTY(execPodShCmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func ConnectToNode() error {
	nodesSelector := NewNodesSelector()
	if err := nodesSelector.PrepareOptions(); err != nil {
		return err
	}

	if len(nodesSelector.nodesInfo) == 0 {
		msg := fmt.Sprintf("\n  No node(s) could be found.")
		logger.Info(msg)
		return nil
	}

	if nodeInfo, err := nodesSelector.SelectNode(); err != nil {
		return err
	} else {
		execNodeBashCmd := fmt.Sprintf("docker exec -it %v %v", nodeInfo.Name, shell.BASH)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + execNodeBashCmd + "\n")
		}

		logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v)", nodeInfo.Name, shell.BASH))
		if err := common.ShellExec.ExecuteTTY(execNodeBashCmd); err != nil {

			// Fallback to /bin/sh if /bin/bash is not available
			execNodeShCmd := fmt.Sprintf("docker exec -it %v %v", nodeInfo.Name, shell.SH)
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + execNodeShCmd + "\n")
			}

			logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v)", nodeInfo.Name, shell.SH))
			if err := common.ShellExec.ExecuteTTY(execNodeShCmd); err != nil {
				return err
			}
		}
	}

	return nil
}
