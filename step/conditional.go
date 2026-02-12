package step

import (
	"encoding/json"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type IfHandler struct{}

func NewIfHandler() *IfHandler {
	return &IfHandler{}
}

func (h *IfHandler) Execute(instance *domain.WorkflowInstance, step domain.StepDefinition) (*runtime.ExecutionResult, error) {
	var spec domain.IfSpec
	if err := json.Unmarshal(step.Spec, &spec); err != nil {
		return nil, err
	}
	result, err := spec.Expression.Evaluate(instance.Variables)
	if err != nil {
		return nil, err
	}
	if result {
		return &runtime.ExecutionResult{
			NextStepID: spec.TrueNext,
		}, nil
	}
	return &runtime.ExecutionResult{
		NextStepID: spec.FalseNext,
	}, nil
}
