package step

import (
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type WaitHandler struct{}

func (h *WaitHandler) Execute(ctx *runtime.Context, step *domain.Step) (executionResult *runtime.ExecutionResult, err error) {
	duration := step.Config["duration_ms"].(int)
	wakeUp := time.Now().Add(time.Duration(duration) * time.Millisecond)
	return &runtime.ExecutionResult{
		Status:   runtime.WAITING,
		NextStep: step.Config["next"].(string),
		WakeUpAt: &wakeUp,
	}, nil
}

func (h *WaitHandler) Validate(step *domain.Step) error {
	return nil
}
