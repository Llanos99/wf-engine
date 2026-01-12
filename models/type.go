package models

const (
	StepIf     StepType = "if-else"
	StepAction StepType = "action"
	StepWait   StepType = "wait"
)

type StepType string

func (s StepType) IsValid() bool {
	switch s {
	case StepIf:
		return true
	case StepAction:
		return true
	default:
		return false
	}
}
