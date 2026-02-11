package step

import (
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type ActionHandler struct{}

func (h *ActionHandler) Execute(ctx *runtime.Context, step *domain.Step) (executionResult *runtime.ExecutionResult, err error) {
	return &runtime.ExecutionResult{
		Status:   runtime.COMPLETED,
		NextStep: step.Config["next"].(string),
	}, nil
}

func (h *ActionHandler) Validate(step *domain.Step) error {
	return nil
}
