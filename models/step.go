package models

type Step struct {
	ID      string
	Name    string
	Type    StepType
	Execute func(*Context) error
	NextID  string
}
