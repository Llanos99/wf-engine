package models

import "fmt"

type Executor struct{}

func (e *Executor) Run(wf *Workflow, ctx *Context) error {
	current := wf.StartAt
	for {
		step := wf.FindStep(current)
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
