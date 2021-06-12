package models

type Instructions struct {
	Items       []*PromptItem `yaml:"promptItems"`
	AutoRun     []string      `yaml:"autoRun"`
	AutoCleanup []string      `yaml:"autoCleanup"`
}
