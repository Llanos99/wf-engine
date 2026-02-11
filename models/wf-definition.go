package models

type WorkflowDefinition struct {
	ID    string             `yaml:"id" json:"id"`
	Start string             `yaml:"start" json:"start"`
	Steps map[string]StepDef `yaml:"steps" json:"steps"`
}

type StepDef struct {
	Type       StepType               `yaml:"type" json:"type"`
	Next       string                 `yaml:"next,omitempty" json:"next,omitempty"`
	TrueNext   string                 `yaml:"true_next,omitempty" json:"true_next,omitempty"`
	FalseNext  string                 `yaml:"false_next,omitempty" json:"false_next,omitempty"`
	Duration   string                 `yaml:"duration,omitempty" json:"duration,omitempty"`
	Action     string                 `yaml:"action,omitempty" json:"action,omitempty"`
	Parameters map[string]interface{} `yaml:"params,omitempty" json:"params,omitempty"`
}
