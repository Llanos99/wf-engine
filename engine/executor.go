package engine

import (
	"errors"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/step"
)

const MAX_LOOP_LIMIT = 20

type Executor struct {
	registry *step.Registry
}

func NewExecutor(registry *step.Registry) *Executor {
	return &Executor{registry: registry}
}

func (e *Executor) Execute(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance) error {
	// While instance.status == running
	for instance.Status == domain.StatusRunning {
		stepDef, ok := def.Steps[instance.CurrentStep]
		if !ok {
			instance.SetStatus(domain.StatusFailed)
			return errors.New("step not found: " + instance.CurrentStep)
		}
		handler := e.registry.Get(stepDef.Type)
		if handler == nil {
			instance.SetStatus(domain.StatusFailed)
			return errors.New("no handler for step type: " + string(stepDef.Type))
		}
		result, err := handler.Execute(instance, stepDef)
		if err != nil {
			instance.SetStatus(domain.StatusFailed)
			return err
		}
		if result.Error != nil {
			instance.SetStatus(domain.StatusFailed)
			return result.Error
		}
		if result.Finished {
			instance.SetStatus(domain.StatusFinished)
			return nil
		}
		if result.WaitUntil != nil {
			instance.SetStatus(domain.StatusWaiting)
			return nil
		}
		if result.NextStepID == "" {
			instance.SetStatus(domain.StatusFinished)
			return nil
		}
		instance.MoveTo(result.NextStepID)
	}
	return nil
}
