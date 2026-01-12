package models

import "time"

type ExecutionStatus int

const (
	COMPLETED ExecutionStatus = iota
	WAITING
	FAILED
)

type ExecutionResult struct {
	Status   ExecutionStatus
	NextStep string
	WakeUpAt *time.Time
}
