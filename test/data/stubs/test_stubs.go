package stubs

import (
	"github.com/ZachiNachshon/anchor/models"
)

const (
	App1Name             = "app-1"
	App1DirPath          = "/app-1"
	App1InstructionsPath = "/app-1/instructions.yaml"

	App2Name             = "app-2"
	App2DirPath          = "/app-2"
	App2InstructionsPath = "/app-2/instructions.yaml"

	App1Action1Id         = "app-1-a"
	App1Action1Title      = "app-1-a title"
	App1Action1ScriptFile = "/path/to/app-1-a"

	App1Action2Id         = "app-1-b"
	App1Action2Title      = "app-1-b title"
	App1Action2ScriptFile = "/path/to/app-1-b"

	App2Action1Id    = "app-2-a"
	App2Action1Title = "app-2-a title"
	App2Action1File  = "/path/to/app-2-a"

	App2Action2Id    = "app-2-b"
	App2Action2Title = "app-2-b title"
	App2Action2File  = "/path/to/app-2-b"

	App1Workflow1Id               = "app-1-w"
	App1Workflow1Description      = "app-1-w description"
	App1Workflow1TolerateFailures = false

	App2Workflow1Id               = "app-2-w"
	App2Workflow1Description      = "app-2-w description"
	App2Workflow1TolerateFailures = false
)

func GenerateApplicationTestData() []*models.ApplicationInfo {
	appDirs := make([]*models.ApplicationInfo, 2)
	appDirs[0] = &app1
	appDirs[1] = &app2
	return appDirs
}

func GenerateInstructionsTestData() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions:   []*models.Action{&app1Action1, &app1Action2},
			Workflows: []*models.Workflow{&app1Workflow1},
		},
	}
}

func GenerateInstructionsTestDataWithoutWorkflows() *models.InstructionsRoot {
	return &models.InstructionsRoot{
		Instructions: &models.Instructions{
			Actions: []*models.Action{&app1Action1, &app1Action2},
		},
	}
}

func GetAppByName(appsArr []*models.ApplicationInfo, name string) *models.ApplicationInfo {
	for _, v := range appsArr {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func GetInstructionActionById(instructions *models.Instructions, id string) *models.Action {
	for _, v := range instructions.Actions {
		if v.Id == id {
			return v
		}
	}
	return nil
}

var app1 = models.ApplicationInfo{
	Name:             App1Name,
	DirPath:          App1DirPath,
	InstructionsPath: App1InstructionsPath,
}

var app1Action1 = models.Action{
	Id:         App1Action1Id,
	Title:      App1Action1Title,
	ScriptFile: App1Action1ScriptFile,
}

var app1Action2 = models.Action{
	Id:         App1Action2Id,
	Title:      App1Action2Title,
	ScriptFile: App1Action2ScriptFile,
}

var app1Workflow1 = models.Workflow{
	Id:               App1Workflow1Id,
	Description:      App1Workflow1Description,
	TolerateFailures: App1Workflow1TolerateFailures,
	ActionIds:        []string{App1Action1Id, App1Action2Id},
}

var app2 = models.ApplicationInfo{
	Name:             App2Name,
	DirPath:          App2DirPath,
	InstructionsPath: App2InstructionsPath,
}

var app2Action1 = models.Action{
	Id:         App2Action1Id,
	Title:      App2Action1Title,
	ScriptFile: App2Action1File,
}

var app2Action2 = models.Action{
	Id:         App2Action2Id,
	Title:      App2Action2Title,
	ScriptFile: App2Action2File,
}

var app2Workflow1 = models.Workflow{
	Id:               App2Workflow1Id,
	Description:      App2Workflow1Description,
	TolerateFailures: App2Workflow1TolerateFailures,
	ActionIds:        []string{App2Action1Id, App2Action2Id},
}
