package engine

import (
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/models"
	"github.com/Llanos99/wf-engine/strategy"
)

type Executor struct{}

func (e *Executor) Run(wf *models.Workflow, ctx *models.Context) error {
	if wf.Validate() != nil {
		return errors.New("Workflow not valid")
	}
	current := wf.StartAt
	for {
		step := wf.FindStepByID(current)
		if !step.Type.IsValid() {
			return fmt.Errorf("Step %s not a valid type", current)
		}
		if step == nil {
			return fmt.Errorf("step %s not found", current)
		}
		handler, ok := strategy.StepHandlers[step.Type]
		if !ok {
			return fmt.Errorf("No handler for step type %s", step.Type)
		}
		next, err := handler.Execute(ctx, step)
		if err != nil {
			return err
		}
		if next == "" {
			return nil
		}
		current = next
	}
}
