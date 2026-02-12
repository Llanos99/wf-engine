package step

import (
	"encoding/json"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type WaitHandler struct{}

func NewWaitHandler() *WaitHandler {
	return &WaitHandler{}
}

func (h *WaitHandler) Execute(instance *domain.WorkflowInstance, step domain.StepDefinition) (*runtime.ExecutionResult, error) {
	var spec domain.WaitSpec
	if err := json.Unmarshal(step.Spec, &spec); err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(spec.Duration)
	if err != nil {
		return nil, err
	}
	waitUntil := time.Now().Add(duration)
	return &runtime.ExecutionResult{
		WaitUntil:  &waitUntil,
		NextStepID: spec.Next,
	}, nil
}
