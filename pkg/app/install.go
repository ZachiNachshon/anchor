package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/banner"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/manifoldco/promptui"
)

func StartApplicationInstallFlow(ctx common.Context) error {
	o, err := orchestrator.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	s, err := shell.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	in, err := input.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	b, err := banner.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	b.PrintAnchor()
	promptErr := runApplicationSelectionFlow(o, s, in)
	if promptErr != nil {
		return managePromptError(promptErr)
	}
	return nil
}

func runApplicationSelectionFlow(o orchestrator.Orchestrator, s shell.Shell, in input.UserInput) *errors.PromptError {
	if app, promptErr := o.OrchestrateApplicationSelection(); promptErr != nil {
		return promptErr
	} else if app.Name == prompter.CancelButtonName {
		return nil
	} else {
		if instructionItem, promptErr := runInstructionSelectionFlow(app, o, s, in); promptErr != nil {
			if promptErr.Code() == errors.InstructionMissingError {
				return runApplicationSelectionFlow(o, s, in)
			}
			return promptErr
		} else if instructionItem.Id == prompter.BackButtonName {
			return runApplicationSelectionFlow(o, s, in)
		}
		return nil
	}
}

func runInstructionSelectionFlow(
	app *models.ApplicationInfo,
	o orchestrator.Orchestrator,
	s shell.Shell,
	in input.UserInput) (*models.InstructionItem, *errors.PromptError) {

	if instructionItem, promptErr := o.OrchestrateInstructionSelection(app); promptErr != nil {
		return nil, promptErr
	} else if instructionItem.Id == prompter.BackButtonName {
		logger.Debugf("Selected to go back from instruction menu. id: %v", instructionItem.Id)
		return instructionItem, nil
	} else {
		logger.Debugf("Selected instruction to run. id: %v", instructionItem.Id)
		if _, promptErr := runInstructionExecutionFlow(instructionItem, o, s, in); promptErr != nil {
			return nil, promptErr
		}
		if inputErr := in.PressAnyKeyToContinue(); inputErr != nil {
			logger.Debugf("Failed to prompt user to press any key after instruction run")
			return nil, errors.New(inputErr)
		} else {
			return runInstructionSelectionFlow(app, o, s, in)
		}
	}
}

func runInstructionExecutionFlow(
	item *models.InstructionItem,
	o orchestrator.Orchestrator,
	s shell.Shell,
	in input.UserInput) (*models.InstructionItem, *errors.PromptError) {

	if shouldRun, promptError := o.AskBeforeRunningInstruction(item, in); promptError != nil {
		logger.Debugf("failed to ask before running an instruction. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		o.RunInstruction(item, s)
	}
	return item, nil
}

func managePromptError(promptErr *errors.PromptError) error {
	err := promptErr.GoError()
	if err == nil {
		logger.Debug("Prompt error returned but does not contain an inner Go error")
		return err
	}
	if err == promptui.ErrInterrupt {
		logger.Debug("exit due to keyboard interrupt")
		return nil
	} else {
		logger.Debug(err.Error())
		return err
	}
}
