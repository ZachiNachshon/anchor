package run

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/runner"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"strings"
)

type DynamicRunFunc func(ctx common.Context, o *runOrchestrator, identifier string) error

var DynamicRun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
	err := o.runner.PrepareFunc(o.runner, ctx)
	if err != nil {
		return err
	}
	err = o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	cmdFolderItem, err := o.extractCommandFolderItemFunc(o, o.commandFolderName, o.commandFolderItemName)
	if err != nil {
		return err
	}
	return o.activeRunFunc(o, ctx, cmdFolderItem, identifier)
}

type runOrchestrator struct {
	verboseFlag bool

	commandFolderName     string
	commandFolderItemName string

	runner *runner.ActionRunnerOrchestrator

	l     locator.Locator
	prntr printer.Printer

	prepareFunc                  func(o *runOrchestrator, ctx common.Context) error
	extractCommandFolderItemFunc func(o *runOrchestrator, commandFolderName string, commandFolderItemName string) (*models.CommandFolderItemInfo, error)
	activeRunFunc                func(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, identifier string) error
	runActionFunc                func(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, identifier string) error
	runWorkflowFunc              func(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, identifier string) error
}

func NewOrchestrator(
	runner *runner.ActionRunnerOrchestrator,
	commandFolderName string,
	commandItemName string) *runOrchestrator {

	return &runOrchestrator{
		commandFolderName:     commandFolderName,
		commandFolderItemName: commandItemName,
		verboseFlag:           false,
		runner:                runner,

		activeRunFunc: func(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, identifier string) error {
			// no-op, shouldn't happen since action/workflow must be used as they are mutually exclusive mandatory flags
			// check tests for additional info
			return nil
		},
		prepareFunc:                  prepare,
		runActionFunc:                runAction,
		runWorkflowFunc:              runWorkflow,
		extractCommandFolderItemFunc: extractCommandFolderItem,
	}
}

func prepare(o *runOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		o.l = resolved.(locator.Locator)
	}
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func extractCommandFolderItem(
	o *runOrchestrator,
	commandFolderName string,
	commandFolderItemName string) (*models.CommandFolderItemInfo, error) {

	folderItems := o.l.CommandFolderItems(commandFolderName)
	if folderItems == nil {
		o.prntr.PrintMissingCommand(commandFolderName)
		return nil, fmt.Errorf("cannot identify dynamic command folder: %s", commandFolderName)
	}

	var commandFolderItemToRun *models.CommandFolderItemInfo = nil
	for _, item := range folderItems {
		if strings.EqualFold(item.Name, commandFolderItemName) {
			commandFolderItemToRun = item
			break
		}
	}

	if commandFolderItemToRun == nil {
		return nil, fmt.Errorf("cannot identify dynamic command folder item: %s", o.commandFolderItemName)
	}
	return commandFolderItemToRun, nil
}

// Usage:
//   - anchor COMMAND_NAME run INSTRUCTION_NAME --action=ACTION_ID
func runAction(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, actionId string) error {
	if instructionsRoot, promptError := o.runner.ExtractInstructionsFunc(o.runner, cmdFolderItem, ctx.AnchorFilesPath()); promptError != nil {
		return promptError.GoError()
	} else {
		action := models.GetInstructionActionById(instructionsRoot.Instructions, actionId)
		if action == nil {
			return fmt.Errorf("%sCannot identify action by id: %s%s", colors.Red, actionId, colors.Reset)
		}

		if promptErr := o.runner.RunInstructionActionFunc(o.runner, action); promptErr != nil {
			return promptErr.GoError()
		}
	}
	return nil
}

// Usage:
//   - anchor COMMAND_NAME run INSTRUCTION_NAME --workflow=WORKFLOW_ID
func runWorkflow(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, workflowId string) error {
	if instructionsRoot, promptError := o.runner.ExtractInstructionsFunc(o.runner, cmdFolderItem, ctx.AnchorFilesPath()); promptError != nil {
		return promptError.GoError()
	} else {
		workflow := models.GetInstructionWorkflowById(instructionsRoot.Instructions, workflowId)
		if workflow == nil {
			return fmt.Errorf("%sCannot identify workflow by id: %s%s", colors.Red, workflowId, colors.Reset)
		}

		actions := instructionsRoot.Instructions.Actions
		if promptErr := o.runner.RunInstructionWorkflowFunc(o.runner, workflow, actions); promptErr != nil {
			return promptErr.GoError()
		}
	}
	return nil
}
