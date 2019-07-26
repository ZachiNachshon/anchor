package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/input"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

type connectCmd struct {
	cobraCmd *cobra.Command
	opts     ConnectCmdOptions
}

type ConnectCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewConnectCmd(opts *common.CmdRootOptions) *connectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a kubernetes pod by name",
		Long:  `Connect to a kubernetes pod by name`,
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Connect to Kubernetes Pod")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {

				_ = loadKubeConfig()

				var namespace = "anchor"
				if len(args) == 2 {
					namespace = args[1]
				}
				if err := connectToPod(args[0], namespace); err != nil {
					logger.Fatal(err.Error())
				}
			}

			logger.PrintCompletion()
		},
	}

	var connectCmd = new(connectCmd)
	connectCmd.cobraCmd = cobraCmd
	connectCmd.opts.CmdRootOptions = opts

	if err := connectCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return connectCmd
}

func connectToPod(name string, namespace string) error {
	podsParser := NewPodsSelector()
	if err := podsParser.PrepareOptions(name, namespace); err != nil {
		return err
	}

	if len(podsParser.podsInfo) == 0 {
		logger.Infof("\n  No %v pods could be found.", name)
		return nil
	}

	return podsParser.SelectPod(name, namespace)
}

func (cmd *connectCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *connectCmd) initFlags() error {
	return nil
}

type podsSelector struct {
	podColumns []*podColumn
	podsInfo   map[int]*podInfo
}

type podInfo struct {
	Name      string
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
	podsColumns = append(podsColumns, &podColumn{Name: "STATUS", Locator: ":.status.phase"})
	podsColumns = append(podsColumns, &podColumn{Name: "VERSION", Locator: ":.metadata.labels.version"})
	podsColumns = append(podsColumns, &podColumn{Name: "START_TIME", Locator: ":.status.startTime"})

	p := &podsSelector{
		podColumns: podsColumns,
	}
	return p
}

func (p *podsSelector) PrepareOptions(name string, namespace string) error {
	if podsInfo, err := p.createPodsInfo(name, namespace); err != nil {
		return err
	} else {
		p.podsInfo = podsInfo
		return nil
	}
}

func (p *podsSelector) SelectPod(name string, namespace string) error {
	if p.podsInfo == nil {
		return errors.Errorf("Something went wrong, pods selector must be prepared prior to pod(s) selection")
	}

	p.printPodsInfo()
	numericInput := input.NewNumericInput()
	maxIdx := len(p.podsInfo)
	success := false
	var rowNum = -1
	for success == false {
		if row, err := numericInput.WaitForInput(); err != nil {
			return err
		} else if row < 1 || row > maxIdx {
			logger.Infof("Selection range can be between [1] and %v\n", maxIdx)
		} else {
			rowNum = row
			success = true
		}
	}

	podInfo := p.podsInfo[rowNum]
	execPodCmd := fmt.Sprintf("kubectl exec -it %v /bin/bash -n %v", podInfo.Name, common.GlobalOptions.KindClusterName)

	if common.GlobalOptions.Verbose {
		logger.Info("\n" + execPodCmd + "\n")
	}

	logger.Infof("\n  Connecting to %v...\n", podInfo.Name)
	if err := common.ShellExec.ExecuteTTY(execPodCmd); err != nil {
		return err
	}
	return nil
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
				Status:    podInfoArr[1],
				Version:   podInfoArr[2],
				StartTime: podInfoArr[3],
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

	lineFormat := "| %v | %-40v %-10v %-10v %-10v\n"
	header := fmt.Sprintf(lineFormat, "#", "NAME", "STATUS", "VERSION", "START_TIME")

	table += header
	for k, v := range p.podsInfo {
		l := fmt.Sprintf(lineFormat, k, v.Name, v.Status, v.Version, v.StartTime)
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
