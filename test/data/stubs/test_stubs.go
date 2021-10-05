package stubs

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

const (
	// Anchor Folder 1
	AnchorFolder1Name         = "app"
	AnchorFolder1DirPath      = "/anchorfiles/app"
	AnchorFolder1Description  = "folder1 description"
	AnchorFolder1CommandUse   = "app"
	AnchorFolder1CommandShort = "Application commands"

	// Anchor Folder 1: Item 1
	AnchorFolder1Item1Name             = "app-1"
	AnchorFolder1Item1DirPath          = "/anchorfiles/app/app-1"
	AnchorFolder1Item1InstructionsPath = "/anchorfiles/app/app-1/instructions.yaml"

	// Anchor Folder 1: Item 1: Action 1
	AnchorFolder1Item1Action1Id          = "app-1-action1"
	AnchorFolder1Item1Action1Title       = "app-1-action1 title"
	AnchorFolder1Item1Action1Description = "app-1-action1 description"
	AnchorFolder1Item1Action1ScriptFile  = "/anchorfiles/app/app-1/script-action1.sh"

	// Anchor Folder 1: Item 1: Action 2
	AnchorFolder1Item1Action2Id          = "app-1-action2"
	AnchorFolder1Item1Action2Title       = "app-1-action2 title"
	AnchorFolder1Item1Action2Description = "app-1-action2 description"
	AnchorFolder1Item1Action2ScriptFile  = "/anchorfiles/app/app-1/script-action2.sh"

	// Anchor Folder 1: Item 1: Workflow 1
	AnchorFolder1Item1Workflow1Id               = "app-1-workflow1"
	AnchorFolder1Item1Workflow1Description      = "app-1-workflow1 description"
	AnchorFolder1Item1Workflow1TolerateFailures = false

	// Anchor Folder 1: Item 2
	AnchorFolder1Item2Name             = "app-2"
	AnchorFolder1Item2DirPath          = "/anchorfiles/app/app-2"
	AnchorFolder2Description           = "folder2 description"
	AnchorFolder1Item2InstructionsPath = "/anchorfiles/app/app-2/instructions.yaml"

	// Anchor Folder 1: Item 2: Action 1
	AnchorFolder1Item2Action1Id          = "app-2-action1"
	AnchorFolder1Item2Action1Title       = "app-2-action1 title"
	AnchorFolder1Item2Action1Description = "app-2-action1 description"
	AnchorFolder1Item2Action1ScriptFile  = "/anchorfiles/app/app-2/script-action1.sh"

	// Anchor Folder 2
	AnchorFolder2Name         = "k8s"
	AnchorFolder2DirPath      = "/anchorfiles/k8s"
	AnchorFolder2CommandUse   = "k8s"
	AnchorFolder2CommandShort = "Kubernetes distributions"

	// Anchor Folder 2: Item 1
	AnchorFolder2Item1Name             = "k8s-1"
	AnchorFolder2Item1DirPath          = "/anchorfiles/k8s/k8s-1"
	AnchorFolder2Item1InstructionsPath = "/anchorfiles/k8s/k8s-1/instructions.yaml"

	// Anchor Folder 2: Item 1: Action 1
	AnchorFolder2Item1Action1Id          = "k8s-1-action1"
	AnchorFolder2Item1Action1Title       = "k8s-1-action1 title"
	AnchorFolder2Item1Action1Description = "k8s-1-action1 description"
	AnchorFolder2Item1Action1Script      = "echo hello k8s-1-action1"

	// Anchor Folder 2: Item 1: Action 2
	AnchorFolder2Item1Action2Id          = "k8s-1-action2"
	AnchorFolder2Item1Action2Title       = "k8s-1-action2 title"
	AnchorFolder2Item1Action2Description = "k8s-1-action2 description"
	AnchorFolder2Item1Action2Script      = "echo hello k8s-1-action2"

	// Anchor Folder 2: Item 1: Workflow 1
	AnchorFolder2Item1Workflow1Id               = "k8s-1-workflow1"
	AnchorFolder2Item1Workflow1Title            = "k8s-1-workflow1 title"
	AnchorFolder2Item1Workflow1Description      = "k8s-1-workflow1 description"
	AnchorFolder2Item1Workflow1TolerateFailures = false
)

func GenerateAnchorFolderInfoTestData() []*models.AnchorFolderInfo {
	anchorFolders := make([]*models.AnchorFolderInfo, 2)
	anchorFolders[0] = anchorFolder1()
	anchorFolders[1] = anchorFolder2()
	return anchorFolders
}

func GenerateAnchorFolderItemsInfoTestData() []*models.AnchorFolderItemInfo {
	anchorFolderItems := make([]*models.AnchorFolderItemInfo, 2)
	anchorFolderItems[0] = anchorFolder1Item1()
	anchorFolderItems[1] = anchorFolder1Item2()
	return anchorFolderItems
}

func GenerateInstructionsTestData() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions:   []*models.Action{anchorFolder1Item1Action1(), anchorFolder1Item1Action2()},
			Workflows: []*models.Workflow{anchorFolder1Item1Workflow1(), app2Workflow1()},
		},
	}
}

func GenerateInstructionsTestDataWithoutWorkflows() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions: []*models.Action{anchorFolder1Item1Action1(), anchorFolder1Item1Action2()},
		},
	}
}

func GetAnchorFolderItemByName(appsArr []*models.AnchorFolderItemInfo, name string) *models.AnchorFolderItemInfo {
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

var anchorFolder1 = func() *models.AnchorFolderInfo {
	item1 := anchorFolder1Item1()
	item2 := anchorFolder1Item2()
	return &models.AnchorFolderInfo{
		Name:    AnchorFolder1Name,
		DirPath: AnchorFolder1DirPath,
		Command: &models.AnchorFolderCommand{
			Use:   AnchorFolder1CommandUse,
			Short: AnchorFolder1CommandShort,
		},
		Description: AnchorFolder1Description,
		Items: map[string]*models.AnchorFolderItemInfo{
			item1.Name: item1,
			item2.Name: item2,
		},
	}
}

var anchorFolder1Item1 = func() *models.AnchorFolderItemInfo {
	return &models.AnchorFolderItemInfo{
		Name:             AnchorFolder1Item1Name,
		DirPath:          AnchorFolder1Item1DirPath,
		InstructionsPath: AnchorFolder1Item1InstructionsPath,
	}
}

var anchorFolder1Item1Action1 = func() *models.Action {
	return &models.Action{
		Id:          AnchorFolder1Item1Action1Id,
		Title:       AnchorFolder1Item1Action1Title,
		Description: AnchorFolder1Item1Action1Description,
		ScriptFile:  AnchorFolder1Item1Action1ScriptFile,
	}
}

var anchorFolder1Item1Action2 = func() *models.Action {
	return &models.Action{
		Id:          AnchorFolder1Item1Action2Id,
		Title:       AnchorFolder1Item1Action2Title,
		Description: AnchorFolder1Item1Action2Description,
		ScriptFile:  AnchorFolder1Item1Action2ScriptFile,
	}
}

var anchorFolder1Item1Workflow1 = func() *models.Workflow {
	return &models.Workflow{
		Id:               AnchorFolder1Item1Workflow1Id,
		Description:      AnchorFolder1Item1Workflow1Description,
		TolerateFailures: AnchorFolder1Item1Workflow1TolerateFailures,
		ActionIds:        []string{AnchorFolder1Item1Action1Id, AnchorFolder1Item1Action2Id},
	}
}

var anchorFolder1Item2 = func() *models.AnchorFolderItemInfo {
	return &models.AnchorFolderItemInfo{
		Name:             AnchorFolder1Item2Name,
		DirPath:          AnchorFolder1Item2DirPath,
		InstructionsPath: AnchorFolder1Item2InstructionsPath,
	}
}

var anchorFolder1Item2Action1 = func() *models.Action {
	return &models.Action{
		Id:          AnchorFolder1Item2Action1Id,
		Title:       AnchorFolder1Item2Action1Title,
		Description: AnchorFolder1Item2Action1Description,
		ScriptFile:  AnchorFolder1Item2Action1ScriptFile,
	}
}

var anchorFolder2 = func() *models.AnchorFolderInfo {
	item1 := anchorFolder2Item1()
	return &models.AnchorFolderInfo{
		Name:    AnchorFolder2Name,
		DirPath: AnchorFolder2DirPath,
		Command: &models.AnchorFolderCommand{
			Use:   AnchorFolder2CommandUse,
			Short: AnchorFolder2CommandShort,
		},
		Description: AnchorFolder2Description,
		Items: map[string]*models.AnchorFolderItemInfo{
			item1.Name: item1,
		},
	}
}

var anchorFolder2Item1 = func() *models.AnchorFolderItemInfo {
	return &models.AnchorFolderItemInfo{
		Name:             AnchorFolder2Item1Name,
		DirPath:          AnchorFolder2Item1DirPath,
		InstructionsPath: AnchorFolder2Item1InstructionsPath,
	}
}

var anchorFolder2Item1Action1 = func() *models.Action {
	return &models.Action{
		Id:          AnchorFolder2Item1Action1Id,
		Title:       AnchorFolder2Item1Action1Title,
		Description: AnchorFolder2Item1Action1Description,
		Script:      AnchorFolder2Item1Action1Script,
	}
}

var anchorFolder2Item1Action2 = func() *models.Action {
	return &models.Action{
		Id:          AnchorFolder2Item1Action2Id,
		Title:       AnchorFolder2Item1Action2Title,
		Description: AnchorFolder2Item1Action2Description,
		Script:      AnchorFolder2Item1Action2Script,
	}
}

var app2Workflow1 = func() *models.Workflow {
	return &models.Workflow{
		Id:               AnchorFolder2Item1Workflow1Id,
		Title:            AnchorFolder2Item1Workflow1Title,
		Description:      AnchorFolder2Item1Workflow1Description,
		TolerateFailures: AnchorFolder2Item1Workflow1TolerateFailures,
		ActionIds:        []string{AnchorFolder2Item1Action1Id, AnchorFolder2Item1Action2Id},
	}
}
