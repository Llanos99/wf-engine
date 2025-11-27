package models

import "log"

type Context struct {
	WorkflowID  string
	InstanceID  string
	Data        map[string]interface{}
	StepResults map[string]interface{}
	Logger      *log.Logger
}
