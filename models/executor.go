package models

import (
	"errors"
	"fmt"
)

type Executor struct{}

func (e *Executor) Run(wf *Workflow, ctx *Context) error {
	if wf.Validate() != nil {
		return errors.New("Workflow not valid")
	}
	current := wf.StartAt
	for {
		step := wf.FindStepByID(current)
		if step == nil {
			return fmt.Errorf("step %s not found", current)
		}
		err := step.Execute(ctx)
		if err != nil {
			return err
		}
		if step.NextID == "" {
			return nil
		}
		current = step.NextID
	}
}
