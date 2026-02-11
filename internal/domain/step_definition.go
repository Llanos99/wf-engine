package domain

import "encoding/json"

type StepType string

const (
	StepTypeIf     StepType = "if-else"
	StepTypeAction StepType = "action"
	StepTypeWait   StepType = "wait"
)

type StepDefinition struct {
	ID   string          `yaml:"id" json:"id"`
	Type StepType        `yaml:"type" json:"type"`
	Spec json.RawMessage `yaml:"spec" json:"spec"`
}
