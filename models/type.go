package models

const (
	StepIf StepType = "if"
)

type StepType string

func (s StepType) IsValid() bool {
	switch s {
	case StepIf:
		return true
	default:
		return false
	}
}
