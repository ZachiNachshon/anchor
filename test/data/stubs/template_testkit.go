package stubs

import (
	"github.com/ZachiNachshon/anchor/models"
)

func GenerateApplicationTestData() []*models.AppContent {
	appDirs := make([]*models.AppContent, 2)
	appDirs[0] = &app1
	appDirs[1] = &app2
	return appDirs
}

func GenerateInstructionsTestData() *models.Instructions {
	return &models.Instructions{
		Items:       []*models.PromptItem{&app1item1, &app1item2},
		AutoRun:     []string{"hello-first-app"},
		AutoCleanup: []string{"goodbye-first-app"},
	}
}

var app1 = models.AppContent{
	Name:             "first-application",
	DirPath:          "/first-app",
	InstructionsPath: "/first-app/instructions.yaml",
}

var app1item1 = models.PromptItem{
	Id:    "hello-world",
	Title: "1st instruction",
	File:  "/path/to/hello-world",
}

var app1item2 = models.PromptItem{
	Id:    "goodbye-world",
	Title: "2nd instruction",
	File:  "/path/to/goodbye-world",
}

var app2 = models.AppContent{
	Name:             "second-application",
	DirPath:          "/second-app",
	InstructionsPath: "/second-app/instructions.yaml",
}

var app2item1 = models.PromptItem{
	Id:    "hello-second-app",
	Title: "This is the 1st item in test",
	File:  "/path/to/hello-second-app",
}

var app2item2 = models.PromptItem{
	Id:    "goodbye-second-app",
	Title: "This is the 2nd item in test",
	File:  "/path/to/goodbye-second-app",
}
