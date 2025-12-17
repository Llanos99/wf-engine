package models

type Step struct {
	ID     string
	Name   string
	Type   StepType
	Config map[string]interface{}
	NextID string // No longer needed, the Handlers decides what the next step is
}
