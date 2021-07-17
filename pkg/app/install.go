package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/banner"
	"github.com/manifoldco/promptui"
)

func StartApplicationInstallFlow(ctx common.Context) error {
	o, err := orchestrator.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	b, err := banner.FromRegistry(ctx.Registry())
	if err != nil {
		return err
	}
	b.PrintAnchorBanner()
	promptErr := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
	if promptErr != nil {
		return managePromptError(promptErr)
	}
	return nil
}

func runApplicationSelectionFlow(o orchestrator.Orchestrator, repoPath string) *errors.PromptError {
	if app, promptErr := o.OrchestrateApplicationSelection(); promptErr != nil {
		return promptErr
	} else if app.Name == prompter.CancelButtonName {
		return nil
	} else {
		if instructionItem, promptErr := runInstructionSelectionFlow(app, o, repoPath); promptErr != nil {
			if promptErr.Code() == errors.InstructionMissingError {
				return runApplicationSelectionFlow(o, repoPath)
			}
			return promptErr
		} else if instructionItem.Id == prompter.BackButtonName {
			return runApplicationSelectionFlow(o, repoPath)
		}
		return nil
	}
}

func runInstructionSelectionFlow(app *models.ApplicationInfo, o orchestrator.Orchestrator, repoPath string) (*models.InstructionItem, *errors.PromptError) {
	if instructionItem, promptErr := o.OrchestrateInstructionSelection(app); promptErr != nil {
		return nil, promptErr
	} else if instructionItem.Id == prompter.BackButtonName {
		logger.Debugf("Selected to go back from instruction menu. id: %v", instructionItem.Id)
		return instructionItem, nil
	} else {
		logger.Debugf("Selected instruction to run. id: %v", instructionItem.Id)
		if _, promptErr := runInstructionExecutionFlow(instructionItem, o, repoPath); promptErr != nil {
			return nil, promptErr
		} else {
			// Clear screen
			return runInstructionSelectionFlow(app, o, repoPath)
		}
	}
}

func runInstructionExecutionFlow(item *models.InstructionItem, o orchestrator.Orchestrator, repoPath string) (*models.InstructionItem, *errors.PromptError) {
	if shouldRun, promptError := o.AskBeforeRunningInstruction(item); promptError != nil {
		logger.Debugf("failed to ask before running an instruction. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.RunInstruction(item, repoPath); promptErr != nil {
			return nil, promptErr
		}
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
