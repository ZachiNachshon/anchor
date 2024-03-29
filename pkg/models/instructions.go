package models

func EmptyGlobals() *Globals {
	return &Globals{}
}

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

func GetInstructionActionById(instructions *Instructions, actionId string) *Action {
	for _, v := range instructions.Actions {
		if v.Id == actionId {
			return v
		}
	}
	return nil
}

func GetInstructionWorkflowById(instructions *Instructions, workflowId string) *Workflow {
	for _, v := range instructions.Workflows {
		if v.Id == workflowId {
			return v
		}
	}
	return nil
}

func GetActionById(actions []*Action, id string) *Action {
	for _, v := range actions {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func GetWorkflowById(workflows []*Workflow, id string) *Workflow {
	for _, v := range workflows {
		if v.Id == id {
			return v
		}
	}
	return nil
}

type InstructionsRoot struct {
	Globals      *Globals      `yaml:"globals"`
	Name         string        `yaml:"name"`
	Instructions *Instructions `yaml:"instructions"`
}

type Globals struct {
	Context string `yaml:"context"`
}

type Instructions struct {
	Actions   []*Action   `yaml:"actions"`
	Workflows []*Workflow `yaml:"workflows"`
}

const (
	ApplicationContext = "application"
	KubernetesContext  = "kubernetes"
)
