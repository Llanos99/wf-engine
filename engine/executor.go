package engine

import (
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
	"github.com/Llanos99/wf-engine/step"
)

type Executor struct{}

const MAX_LOOP_LIMIT = 20

func (e *Executor) Run(wf *domain.Workflow, ctx *runtime.Context) error {
	if wf.Validate() != nil {
		return errors.New("Workflow not valid")
	}
	current := wf.StartAt
	for {
		currStep := wf.FindStepByID(current)

		if currStep == nil {
			return fmt.Errorf("currStep %s not found", current)
		}

		ctx.ExecutionCount[currStep.ID] += 1
		if ctx.ExecutionCount[currStep.ID] > MAX_LOOP_LIMIT {
			return fmt.Errorf("Step %s has exceeded the max executions (%d)", currStep.ID, MAX_LOOP_LIMIT)
		}
		handler, ok := step.StepHandlers[currStep.Type]
		if !ok {
			return fmt.Errorf("No handler for currStep type %s", currStep.Type)
		}

		if handlerIsValid := handler.Validate(currStep); handlerIsValid != nil {
			return fmt.Errorf("Handler for currStep %s is not valid", currStep.ID)
		}

		result, err := handler.Execute(ctx, currStep)

		if err != nil {
			return err
		}

		if result.NextStep == "" {
			return nil
		}

		switch result.Status {
		case runtime.COMPLETED:
			current = result.NextStep
		case runtime.WAITING:
			// Persist state and exit
			return nil
		case runtime.FAILED:
			return fmt.Errorf("Execution of currStep %s failed", currStep.ID)
		}
	}
}
