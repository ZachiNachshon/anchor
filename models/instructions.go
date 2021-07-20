package models

func EmptyInstructionsRoot() *InstructionsRoot {
	return &InstructionsRoot{
		Instructions: &Instructions{
			Actions:   make([]*Action, 0, 0),
			Workflows: make([]*Workflow, 0, 0),
		},
	}
}

type InstructionsRoot struct {
	Instructions *Instructions `yaml:"instructions"`
}

type Instructions struct {
	Actions   []*Action   `yaml:"actions"`
	Workflows []*Workflow `yaml:"workflows"`
}
