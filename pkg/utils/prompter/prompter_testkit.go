package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
)

func GenerateApplicationTestData() []*locator.AppContent {
	appDirs := make([]*locator.AppContent, 2)
	appDirs[0] = &app1
	appDirs[1] = &app2
	return appDirs
}

func GenerateInstructionsTestData() *parser.Instructions {
	return &parser.Instructions{
		Items:       []*parser.PromptItem{&app1item1, &app1item2},
		AutoRun:     []string{"hello-first-app"},
		AutoCleanup: []string{"goodbye-first-app"},
	}
}

var app1 = locator.AppContent{
	Name:             "first-application",
	DirPath:          "/first-app",
	InstructionsPath: "/first-app/instructions.yaml",
}

var app1item1 = parser.PromptItem{
	Id:    "hello-first-app",
	Title: "This is the 1st item in test",
	File:  "/path/to/hello-first-app",
}

var app1item2 = parser.PromptItem{
	Id:    "goodbye-first-app",
	Title: "This is the 2nd item in test",
	File:  "/path/to/goodbye-first-app",
}

var app2 = locator.AppContent{
	Name:             "second-application",
	DirPath:          "/second-app",
	InstructionsPath: "/second-app/instructions.yaml",
}

var app2item1 = parser.PromptItem{
	Id:    "hello-second-app",
	Title: "This is the 1st item in test",
	File:  "/path/to/hello-second-app",
}

var app2item2 = parser.PromptItem{
	Id:    "goodbye-second-app",
	Title: "This is the 2nd item in test",
	File:  "/path/to/goodbye-second-app",
}
