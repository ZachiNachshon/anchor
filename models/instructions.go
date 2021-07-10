package models

type Instructions struct {
	Items       []*InstructionItem `yaml:"promptItems"`
	AutoRun     []string           `yaml:"autoRun"`
	AutoCleanup []string           `yaml:"autoCleanup"`
}
