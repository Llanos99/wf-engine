package models

import "log"

type Context struct {
	WorkflowID     string
	InstanceID     string
	Variables      *Variables
	StepResults    map[string]interface{}
	Logger         *log.Logger
	ExecutionCount map[string]int
}
