package stubs

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

const (
	// Anchor Folder 1
	CommandFolder1Name         = "app"
	CommandFolder1DirPath      = "/anchorfiles/app"
	CommandFolder1Description  = "folder1 description"
	CommandFolder1CommandUse   = "app"
	CommandFolder1CommandShort = "Application commands"

	// Anchor Folder 1: Item 1
	CommandFolder1Item1Name             = "app-1"
	CommandFolder1Item1DirPath          = "/anchorfiles/app/app-1"
	CommandFolder1Item1InstructionsPath = "/anchorfiles/app/app-1/instructions.yaml"

	// Anchor Folder 1: Item 1: Action 1
	CommandFolder1Item1Action1Id          = "app-1-action1"
	CommandFolder1Item1Action1Title       = "app-1-action1 title"
	CommandFolder1Item1Action1Description = "app-1-action1 description"
	CommandFolder1Item1Action1ScriptFile  = "/anchorfiles/app/app-1/script-action1.sh"

	// Anchor Folder 1: Item 1: Action 2
	CommandFolder1Item1Action2Id          = "app-1-action2"
	CommandFolder1Item1Action2Title       = "app-1-action2 title"
	CommandFolder1Item1Action2Description = "app-1-action2 description"
	CommandFolder1Item1Action2ScriptFile  = "/anchorfiles/app/app-1/script-action2.sh"

	// Anchor Folder 1: Item 1: Workflow 1
	CommandFolder1Item1Workflow1Id               = "app-1-workflow1"
	CommandFolder1Item1Workflow1Description      = "app-1-workflow1 description"
	CommandFolder1Item1Workflow1TolerateFailures = false

	// Anchor Folder 1: Item 2
	CommandFolder1Item2Name             = "app-2"
	CommandFolder1Item2DirPath          = "/anchorfiles/app/app-2"
	CommandFolder2Description           = "folder2 description"
	CommandFolder1Item2InstructionsPath = "/anchorfiles/app/app-2/instructions.yaml"

	// Anchor Folder 1: Item 2: Action 1
	CommandFolder1Item2Action1Id          = "app-2-action1"
	CommandFolder1Item2Action1Title       = "app-2-action1 title"
	CommandFolder1Item2Action1Description = "app-2-action1 description"
	CommandFolder1Item2Action1ScriptFile  = "/anchorfiles/app/app-2/script-action1.sh"

	// Anchor Folder 2
	CommandFolder2Name         = "k8s"
	CommandFolder2DirPath      = "/anchorfiles/k8s"
	CommandFolder2CommandUse   = "k8s"
	CommandFolder2CommandShort = "Kubernetes distributions"

	// Anchor Folder 2: Item 1
	CommandFolder2Item1Name             = "k8s-1"
	CommandFolder2Item1DirPath          = "/anchorfiles/k8s/k8s-1"
	CommandFolder2Item1InstructionsPath = "/anchorfiles/k8s/k8s-1/instructions.yaml"

	// Anchor Folder 2: Item 1: Action 1
	CommandFolder2Item1Action1Id          = "k8s-1-action1"
	CommandFolder2Item1Action1Title       = "k8s-1-action1 title"
	CommandFolder2Item1Action1Description = "k8s-1-action1 description"
	CommandFolder2Item1Action1Script      = "echo hello k8s-1-action1"

	// Anchor Folder 2: Item 1: Action 2
	CommandFolder2Item1Action2Id          = "k8s-1-action2"
	CommandFolder2Item1Action2Title       = "k8s-1-action2 title"
	CommandFolder2Item1Action2Description = "k8s-1-action2 description"
	CommandFolder2Item1Action2Script      = "echo hello k8s-1-action2"

	// Anchor Folder 2: Item 1: Workflow 1
	CommandFolder2Item1Workflow1Id               = "k8s-1-workflow1"
	CommandFolder2Item1Workflow1Title            = "k8s-1-workflow1 title"
	CommandFolder2Item1Workflow1Description      = "k8s-1-workflow1 description"
	CommandFolder2Item1Workflow1TolerateFailures = false
)

func GenerateCommandFolderInfoTestData() []*models.CommandFolderInfo {
	commandFolders := make([]*models.CommandFolderInfo, 2)
	commandFolders[0] = commandFolder1()
	commandFolders[1] = commandFolder2()
	return commandFolders
}

func GenerateCommandFolderItemsInfoTestData() []*models.CommandFolderItemInfo {
	commandFolderItems := make([]*models.CommandFolderItemInfo, 2)
	commandFolderItems[0] = commandFolder1Item1()
	commandFolderItems[1] = commandFolder1Item2()
	return commandFolderItems
}

func GenerateInstructionsTestData() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions:   []*models.Action{commandFolder1Item1Action1(), commandFolder1Item1Action2()},
			Workflows: []*models.Workflow{commandFolder1Item1Workflow1(), app2Workflow1()},
		},
	}
}

func GenerateInstructionsTestDataWithoutWorkflows() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions: []*models.Action{commandFolder1Item1Action1(), commandFolder1Item1Action2()},
		},
	}
}

func GetCommandFolderItemByName(appsArr []*models.CommandFolderItemInfo, name string) *models.CommandFolderItemInfo {
	for _, v := range appsArr {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func GetInstructionActionById(instructions *models.Instructions, actionId string) *models.Action {
	for _, v := range instructions.Actions {
		if v.Id == actionId {
			return v
		}
	}
	return nil
}

func GetInstructionWorkflowById(instructions *models.Instructions, workflowId string) *models.Workflow {
	for _, v := range instructions.Workflows {
		if v.Id == workflowId {
			return v
		}
	}
	return nil
}

var commandFolder1 = func() *models.CommandFolderInfo {
	item1 := commandFolder1Item1()
	item2 := commandFolder1Item2()
	return &models.CommandFolderInfo{
		Name:    CommandFolder1Name,
		DirPath: CommandFolder1DirPath,
		Command: &models.CommandFolderCommand{
			Use:   CommandFolder1CommandUse,
			Short: CommandFolder1CommandShort,
		},
		Description: CommandFolder1Description,
		Items: map[string]*models.CommandFolderItemInfo{
			item1.Name: item1,
			item2.Name: item2,
		},
	}
}

var commandFolder1Item1 = func() *models.CommandFolderItemInfo {
	return &models.CommandFolderItemInfo{
		Name:             CommandFolder1Item1Name,
		DirPath:          CommandFolder1Item1DirPath,
		InstructionsPath: CommandFolder1Item1InstructionsPath,
	}
}

var commandFolder1Item1Action1 = func() *models.Action {
	return &models.Action{
		Id:          CommandFolder1Item1Action1Id,
		Title:       CommandFolder1Item1Action1Title,
		Description: CommandFolder1Item1Action1Description,
		ScriptFile:  CommandFolder1Item1Action1ScriptFile,
	}
}

var commandFolder1Item1Action2 = func() *models.Action {
	return &models.Action{
		Id:          CommandFolder1Item1Action2Id,
		Title:       CommandFolder1Item1Action2Title,
		Description: CommandFolder1Item1Action2Description,
		ScriptFile:  CommandFolder1Item1Action2ScriptFile,
	}
}

var commandFolder1Item1Workflow1 = func() *models.Workflow {
	return &models.Workflow{
		Id:               CommandFolder1Item1Workflow1Id,
		Description:      CommandFolder1Item1Workflow1Description,
		TolerateFailures: CommandFolder1Item1Workflow1TolerateFailures,
		ActionIds:        []string{CommandFolder1Item1Action1Id, CommandFolder1Item1Action2Id},
	}
}

var commandFolder1Item2 = func() *models.CommandFolderItemInfo {
	return &models.CommandFolderItemInfo{
		Name:             CommandFolder1Item2Name,
		DirPath:          CommandFolder1Item2DirPath,
		InstructionsPath: CommandFolder1Item2InstructionsPath,
	}
}

var commandFolder1Item2Action1 = func() *models.Action {
	return &models.Action{
		Id:          CommandFolder1Item2Action1Id,
		Title:       CommandFolder1Item2Action1Title,
		Description: CommandFolder1Item2Action1Description,
		ScriptFile:  CommandFolder1Item2Action1ScriptFile,
	}
}

var commandFolder2 = func() *models.CommandFolderInfo {
	item1 := commandFolder2Item1()
	return &models.CommandFolderInfo{
		Name:    CommandFolder2Name,
		DirPath: CommandFolder2DirPath,
		Command: &models.CommandFolderCommand{
			Use:   CommandFolder2CommandUse,
			Short: CommandFolder2CommandShort,
		},
		Description: CommandFolder2Description,
		Items: map[string]*models.CommandFolderItemInfo{
			item1.Name: item1,
		},
	}
}

var commandFolder2Item1 = func() *models.CommandFolderItemInfo {
	return &models.CommandFolderItemInfo{
		Name:             CommandFolder2Item1Name,
		DirPath:          CommandFolder2Item1DirPath,
		InstructionsPath: CommandFolder2Item1InstructionsPath,
	}
}

var commandFolder2Item1Action1 = func() *models.Action {
	return &models.Action{
		Id:          CommandFolder2Item1Action1Id,
		Title:       CommandFolder2Item1Action1Title,
		Description: CommandFolder2Item1Action1Description,
		Script:      CommandFolder2Item1Action1Script,
	}
}

var commandFolder2Item1Action2 = func() *models.Action {
	return &models.Action{
		Id:          CommandFolder2Item1Action2Id,
		Title:       CommandFolder2Item1Action2Title,
		Description: CommandFolder2Item1Action2Description,
		Script:      CommandFolder2Item1Action2Script,
	}
}

var app2Workflow1 = func() *models.Workflow {
	return &models.Workflow{
		Id:               CommandFolder2Item1Workflow1Id,
		Title:            CommandFolder2Item1Workflow1Title,
		Description:      CommandFolder2Item1Workflow1Description,
		TolerateFailures: CommandFolder2Item1Workflow1TolerateFailures,
		ActionIds:        []string{CommandFolder2Item1Action1Id, CommandFolder2Item1Action2Id},
	}
}
