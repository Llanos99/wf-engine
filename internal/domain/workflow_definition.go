package domain

type WorkflowDefinition struct {
	ID               string                    `yaml:"id" json:"id"`
	Version          int                       `yaml:"version" json:"version"`
	Name             string                    `yaml:"name" json:"name"`
	Description      string                    `yaml:"description" json:"description"`
	StartAt          string                    `yaml:"start_at" json:"start_at"`
	Steps            map[string]StepDefinition `yaml:"steps" json:"steps"`
	InitialVariables Variables                 `yaml:"initial_variables" json:"initial_variables"`
}
