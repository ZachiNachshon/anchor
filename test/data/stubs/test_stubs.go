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

	App1InstructionsItem1Id    = "app-1-a"
	App1InstructionsItem1Title = "app-1-a title"
	App1InstructionsItem1File  = "/path/to/app-1-a"

	App1InstructionsItem2Id    = "app-1-b"
	App1InstructionsItem2Title = "app-1-b title"
	App1InstructionsItem2File  = "/path/to/app-1-b"

	App2InstructionsItem1Id    = "app-2-a"
	App2InstructionsItem1Title = "app-2-a title"
	App2InstructionsItem1File  = "/path/to/app-2-a"

	App2InstructionsItem2Id    = "app-2-b"
	App2InstructionsItem2Title = "app-2-b title"
	App2InstructionsItem2File  = "/path/to/app-2-b"
)

func GenerateApplicationTestData() []*models.AppContent {
	appDirs := make([]*models.AppContent, 2)
	appDirs[0] = &app1
	appDirs[1] = &app2
	return appDirs
}

func GenerateInstructionsTestData() *models.Instructions {
	return &models.Instructions{
		Items:       []*models.PromptItem{&app1InstructionsItem1, &app1InstructionsItem2},
		AutoRun:     []string{App1InstructionsItem1Id},
		AutoCleanup: []string{App1InstructionsItem2Id},
	}
}

var app1 = models.AppContent{
	Name:             App1Name,
	DirPath:          App1DirPath,
	InstructionsPath: App1InstructionsPath,
}

var app1InstructionsItem1 = models.PromptItem{
	Id:    App1InstructionsItem1Id,
	Title: App1InstructionsItem1Title,
	File:  App1InstructionsItem1File,
}

var app1InstructionsItem2 = models.PromptItem{
	Id:    App1InstructionsItem2Id,
	Title: App1InstructionsItem2Title,
	File:  App1InstructionsItem2File,
}

var app2 = models.AppContent{
	Name:             App2Name,
	DirPath:          App2DirPath,
	InstructionsPath: App2InstructionsPath,
}

var app2InstructionsItem1 = models.PromptItem{
	Id:    App2InstructionsItem1Id,
	Title: App2InstructionsItem1Title,
	File:  App2InstructionsItem1File,
}

var app2InstructionsItem2 = models.PromptItem{
	Id:    App2InstructionsItem2Id,
	Title: App2InstructionsItem2Title,
	File:  App2InstructionsItem2File,
}
