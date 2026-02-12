package step

import (
	"encoding/json"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type ActionHandler struct{}

func NewActionHandler() *ActionHandler {
	return &ActionHandler{}
}

func (h *ActionHandler) Execute(instance *domain.WorkflowInstance, step domain.StepDefinition) (*runtime.ExecutionResult, error) {
	var spec domain.ActionSpec
	if err := json.Unmarshal(step.Spec, &spec); err != nil {
		return nil, err
	}
	switch spec.ActionName {
	case "compute_even":
		nFloat := instance.Variables["n"].(float64)
		n := int(nFloat)
		instance.Variables["even"] = n%2 == 0
		fmt.Println("n =", n)
	case "divide_by_2":
		nFloat := instance.Variables["n"].(float64)
		n := int(nFloat)
		instance.Variables["n"] = float64(n / 2)
	case "multiply_3n_plus_1":
		nFloat := instance.Variables["n"].(float64)
		n := int(nFloat)
		instance.Variables["n"] = float64(3*n + 1)
	case "finish":
		return &runtime.ExecutionResult{
			Finished: true,
		}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", spec.ActionName)
	}
	return &runtime.ExecutionResult{
		NextStepID: spec.Next,
	}, nil
}
