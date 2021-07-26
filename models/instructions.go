package models

func EmptyActions() []*Action {
	return make([]*Action, 0, 0)
}

func EmptyInstructionsRoot() *InstructionsRoot {
	return &InstructionsRoot{
		Instructions: &Instructions{
			Actions:   make([]*Action, 0, 0),
			Workflows: make([]*Workflow, 0, 0),
		},
	}
}

func GetInstructionActionById(actions []*Action, id string) *Action {
	for _, v := range actions {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func GetInstructionWorkflowById(workflows []*Workflow, id string) *Workflow {
	for _, v := range workflows {
		if v.Id == id {
			return v
		}
	}
	return nil
}

type InstructionsRoot struct {
	Instructions *Instructions `yaml:"instructions"`
}

type Instructions struct {
	Actions   []*Action   `yaml:"actions"`
	Workflows []*Workflow `yaml:"workflows"`
}
