package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const statefulLabel = "stateful"

type nodesSelector struct {
	nodeColumns []*nodeColumn
	nodesInfo   map[int]*nodeInfo
}

type nodeInfo struct {
	Name    string
	Status  string
	Roles   string
	Age     string
	Version string
}

type nodeColumn struct {
	Name    string
	Locator string
}

func NewNodesSelector() *nodesSelector {
	var nodesColumns []*nodeColumn

	//nodesColumns = append(nodesColumns, &nodeColumn{Name: "NAME", Locator: ":.metadata.name"})
	//nodesColumns = append(nodesColumns, &nodeColumn{Name: "STATUS", Locator: ":.status.phase"})
	//nodesColumns = append(nodesColumns, &nodeColumn{Name: "VERSION", Locator: ":.metadata.labels.version"})
	//nodesColumns = append(nodesColumns, &nodeColumn{Name: "START_TIME", Locator: ":.status.startTime"})

	n := &nodesSelector{
		nodeColumns: nodesColumns,
	}
	return n
}

func (n *nodesSelector) createNodesInfo() (map[int]*nodeInfo, error) {
	execCmd := n.getSelectableNodesCommand()
	nodesInfo := make(map[int]*nodeInfo)
	if nodesOutput, err := common.ShellExec.ExecuteWithOutput(execCmd); err != nil {
		return nil, err
	} else {
		lines := strings.Split(nodesOutput, "\n")
		for i, line := range lines {
			// Ignore header & empty lines
			if i == 0 || len(line) == 0 {
				continue
			}

			nodeInfoArr := strings.Fields(line)
			info := &nodeInfo{
				Name:    nodeInfoArr[0],
				Status:  nodeInfoArr[1],
				Roles:   nodeInfoArr[2],
				Age:     nodeInfoArr[3],
				Version: nodeInfoArr[4],
			}

			nodesInfo[i] = info
		}
	}
	return nodesInfo, nil
}

func (n *nodesSelector) printPodsInfo() {
	if n.nodesInfo == nil {
		logger.Info("Something went wrong, missing node(s) information")
		return
	}

	table := "\n"

	lineFormat := "| %v | %-25v %-10v %-10v %-10v %-10v\n"
	header := fmt.Sprintf(lineFormat, "#", "NAME", "STATUS", "ROLE", "AGE", "VERSION")

	table += header

	// Keep the numeric ordering
	for i := 1; i <= len(n.nodesInfo); i++ {
		v := n.nodesInfo[i]
		l := fmt.Sprintf(lineFormat, i, v.Name, v.Status, v.Roles, v.Age, v.Version)
		table += l
	}

	logger.Info(table)
}

func (n *nodesSelector) getSelectableNodesCommand() string {
	getNodesFmt := "kubectl get nodes"

	// TODO: Explicitly decide on which columns to collect
	//columnsOpt := n.getColumnsAsKubectlOption()

	getPodsCmd := fmt.Sprintf(getNodesFmt)
	return getPodsCmd
}

func (n *nodesSelector) getColumnsAsKubectlOption() string {
	customColumns := "-o custom-columns="
	del := ""
	for _, c := range n.nodeColumns {
		customColumns += del + c.Name + c.Locator
		del = ","
	}
	return customColumns
}

func (n *nodesSelector) PrepareOptions() error {
	if nodesInfo, err := n.createNodesInfo(); err != nil {
		return err
	} else {
		n.nodesInfo = nodesInfo
		return nil
	}
}

func (n *nodesSelector) SelectNode() (*nodeInfo, error) {
	if n.nodesInfo == nil {
		return nil, errors.Errorf("Something went wrong, nodes selector must be prepared prior to pod(s) selection")
	}

	n.printPodsInfo()
	numericInput := input.NewNumericInput()
	maxIdx := len(n.nodesInfo)
	success := false
	var rowNum = -1
	for success == false {
		if row, err := numericInput.WaitForInput(); err != nil {
			return nil, err
		} else if row < 1 || row > maxIdx {
			logger.Infof("Selection range can be between [1] and %v\n", maxIdx)
		} else {
			rowNum = row
			success = true
		}
	}

	nodeInfo := n.nodesInfo[rowNum]
	return nodeInfo, nil
}

func getNodeNameByLabel(label string) (string, error) {
	getNodeByLblCmd := fmt.Sprintf("kubectl get node --no-headers -l \"%v\" | awk '{print $1}'", label)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + getNodeByLblCmd + "\n")
	}
	if out, err := common.ShellExec.ExecuteWithOutput(getNodeByLblCmd); err != nil {
		return "", err
	} else {
		return strings.Trim(out, "\n"), nil
	}
}

type podsSelector struct {
	podColumns []*podColumn
	podsInfo   map[int]*podInfo
}

type podInfo struct {
	Name      string
	Node      string
	Status    string
	Version   string
	StartTime string
}

type podColumn struct {
	Name    string
	Locator string
}

func NewPodsSelector() *podsSelector {
	var podsColumns []*podColumn

	podsColumns = append(podsColumns, &podColumn{Name: "NAME", Locator: ":.metadata.name"})
	podsColumns = append(podsColumns, &podColumn{Name: "NODE", Locator: ":.spec.nodeName"})
	podsColumns = append(podsColumns, &podColumn{Name: "STATUS", Locator: ":.status.phase"})
	podsColumns = append(podsColumns, &podColumn{Name: "VERSION", Locator: ":.metadata.labels.version"})
	podsColumns = append(podsColumns, &podColumn{Name: "START_TIME", Locator: ":.status.startTime"})

	p := &podsSelector{
		podColumns: podsColumns,
	}
	return p
}

func (p *podsSelector) createPodsInfo(name string, namespace string) (map[int]*podInfo, error) {
	execCmd := p.getSelectablePodsCommand(name, namespace)
	podsInfo := make(map[int]*podInfo)
	if podsOutput, err := common.ShellExec.ExecuteWithOutput(execCmd); err != nil {
		return nil, err
	} else {
		lines := strings.Split(podsOutput, "\n")
		for i, line := range lines {
			// Ignore header & empty lines
			if i == 0 || len(line) == 0 {
				continue
			}

			podInfoArr := strings.Fields(line)
			info := &podInfo{
				Name:      podInfoArr[0],
				Node:      podInfoArr[1],
				Status:    podInfoArr[2],
				Version:   podInfoArr[3],
				StartTime: podInfoArr[4],
			}

			podsInfo[i] = info
		}
	}
	return podsInfo, nil
}

func (p *podsSelector) printPodsInfo() {
	if p.podsInfo == nil {
		logger.Info("Something went wrong, missing pod(s) information")
		return
	}

	table := "\n"

	lineFormat := "| %v | %-35v %-17v %-10v %-10v %-10v\n"
	header := fmt.Sprintf(lineFormat, "#", "NAME", "NODE", "STATUS", "VERSION", "START_TIME")

	table += header

	// Keep the numeric ordering
	for i := 1; i <= len(p.podsInfo); i++ {
		v := p.podsInfo[i]
		l := fmt.Sprintf(lineFormat, i, v.Name, v.Node, v.Status, v.Version, v.StartTime)
		table += l
	}

	logger.Info(table)
}

func (p *podsSelector) getSelectablePodsCommand(name string, namespace string) string {
	getPodsFmt := "kubectl get pods %v %v | grep -Ei \"%v|NAME\""

	namespaceOpt := fmt.Sprintf("-n %v", namespace)
	columnsOpt := p.getColumnsAsKubectlOption()

	getPodsCmd := fmt.Sprintf(getPodsFmt, namespaceOpt, columnsOpt, name)
	return getPodsCmd
}

func (p *podsSelector) getColumnsAsKubectlOption() string {
	customColumns := "-o custom-columns="
	del := ""
	for _, c := range p.podColumns {
		customColumns += del + c.Name + c.Locator
		del = ","
	}
	return customColumns
}

func (p *podsSelector) PrepareOptions(name string, namespace string) error {
	if podsInfo, err := p.createPodsInfo(name, namespace); err != nil {
		return err
	} else {
		p.podsInfo = podsInfo
		return nil
	}
}

func (p *podsSelector) SelectPod(name string, namespace string) (*podInfo, error) {
	if p.podsInfo == nil {
		return nil, errors.Errorf("Something went wrong, pods selector must be prepared prior to pod(s) selection")
	}

	p.printPodsInfo()
	numericInput := input.NewNumericInput()
	maxIdx := len(p.podsInfo)
	success := false
	var rowNum = -1
	for success == false {
		if row, err := numericInput.WaitForInput(); err != nil {
			return nil, err
		} else if row < 1 || row > maxIdx {
			logger.Infof("Selection range can be between [1] and %v\n", maxIdx)
		} else {
			rowNum = row
			success = true
		}
	}

	podInfo := p.podsInfo[rowNum]
	return podInfo, nil
}

func addNodeLabel(node string, namespace string, key string, value string) error {
	logger.PrintCommandHeader(fmt.Sprintf("Added label '%v=%v' to node %v", key, statefulLabel, node))
	labelNodeCmd := fmt.Sprintf("kubectl label node %v %v=%v --overwrite=true -n %v", node, key, value, namespace)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + labelNodeCmd + "\n")
	}
	return common.ShellExec.Execute(labelNodeCmd)
}

func removeNodeLabel(node string, namespace string, key string) error {
	logger.PrintCommandHeader(fmt.Sprintf("Removing label '%v=%v' from node %v", key, statefulLabel, node))
	unlabelNodeCmd := fmt.Sprintf("kubectl label node %v %v- -n %v", node, key, namespace)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + unlabelNodeCmd + "\n")
	}
	return common.ShellExec.Execute(unlabelNodeCmd)
}

func LoadKubeConfig() error {
	// Export k8s configuration
	loadCmd := "kind get kubeconfig-path --name=" + common.GlobalOptions.KindClusterName
	if out, err := common.ShellExec.ExecuteWithOutput(loadCmd); err != nil {
		return err
	} else {
		out = strings.TrimSuffix(out, "\n")
		return os.Setenv("KUBECONFIG", out)
	}
}

func CheckForActiveCluster(name string) (bool, error) {
	getClustersCmd := "kind get clusters"
	if out, err := common.ShellExec.ExecuteWithOutput(getClustersCmd); err != nil {
		return false, err
	} else {
		contains := strings.Contains(out, name)
		return contains, nil
	}
}

func KillKubectlProxy() error {
	return common.ShellExec.Execute(`ps -ef | grep "kubectl proxy" | grep -v grep | awk '{print $2}' | xargs kill -9`)
}

func Prerequisites() bool {
	name := common.GlobalOptions.KindClusterName
	if exists, err := CheckForActiveCluster(name); err != nil {
		logger.Fatal(err.Error())
	} else if !exists {
		logger.Info("\n No active cluster.\n")
		return false
	}

	_ = LoadKubeConfig()

	return true
}

func extractVolumePath(hostPath string) string {
	volumePath := hostPath[strings.Index(hostPath, " "):]
	volumePath = strings.TrimSpace(volumePath)
	volumePath = strings.TrimSuffix(volumePath, "\n")
	return volumePath
}

func mountHostPath(name string, namespace string, hostPath string) error {
	volumePath := extractVolumePath(hostPath)

	// Init node selector and prepare node options
	nodesSelector := NewNodesSelector()
	if err := nodesSelector.PrepareOptions(); err != nil {
		return err
	}

	if len(nodesSelector.nodesInfo) == 0 {
		msg := fmt.Sprintf("  No node(s) could be found.\n")
		logger.Info(msg)
		return errors.Errorf("Failed applying manifest since no nodes could be found")
	}

	if nodeInfo, err := nodesSelector.SelectNode(); err != nil {
		return err
	} else {
		// Label node as a stateful node for deployed content
		if err := addNodeLabel(nodeInfo.Name, namespace, name, statefulLabel); err != nil {
			return err
		} else {
			// Copy hostPath content to <node>/opt/stateful
			logger.PrintCommandHeader(fmt.Sprintf("Copying %v to %v:/opt/stateful/%v", volumePath, nodeInfo.Name, name))
			copyHostPathCmd := fmt.Sprintf("docker cp %v %v:/opt/stateful/%v", volumePath, nodeInfo.Name, name)

			if common.GlobalOptions.Verbose {
				logger.Info("\n" + copyHostPathCmd + "\n")
			}

			if err := common.ShellExec.Execute(copyHostPathCmd); err != nil {
				return err
			}
		}
	}
	return nil
}

func unMountHostPath(name string, namespace string, hostPath string) error {
	volumePath := extractVolumePath(hostPath)

	label := fmt.Sprintf("%v=%v", name, statefulLabel)
	var nodeName string
	var err error
	if nodeName, err = getNodeNameByLabel(label); err != nil {
		return err
	} else if len(nodeName) == 0 {
		msg := fmt.Sprintf("  No node(s) could be found with label %v.\n", label)
		logger.Info(msg)
	} else {

		if err := backupMountPath(nodeName, name, volumePath); err != nil {
			return err
		}

		if err := deleteMountPath(nodeName, name); err != nil {
			return err
		}

		if err := removeNodeLabel(nodeName, namespace, name); err != nil {
			return err
		}
	}
	return nil
}

func backupMountPath(nodeName string, name string, volumePath string) error {
	// Copy content from <node>/opt/stateful/<name> to hostPath
	logger.PrintCommandHeader(fmt.Sprintf("Copying %v:/opt/stateful/%v to %v", nodeName, name, volumePath))
	copyMountPathCmd := fmt.Sprintf("docker cp %v:/opt/stateful/%v %v", nodeName, name, volumePath)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + copyMountPathCmd + "\n")
	}
	if err := common.ShellExec.Execute(copyMountPathCmd); err != nil {
		return err
	}
	return nil
}

func deleteMountPath(nodeName string, name string) error {
	// Delete mounted volume on node
	logger.PrintCommandHeader(fmt.Sprintf("Deleting %v:/opt/stateful/%v", nodeName, name))
	deleteVolumeCmd := fmt.Sprintf("docker exec -t %v /bin/bash 'rm -rf /opt/stateful/%v'", nodeName, name)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + deleteVolumeCmd + "\n")
	}
	if err := common.ShellExec.Execute(deleteVolumeCmd); err != nil {
		return err
	}
	return nil
}
