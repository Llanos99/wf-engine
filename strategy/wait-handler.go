package strategy

import (
	"time"

	"github.com/Llanos99/wf-engine/models"
)

type WaitHandler struct{}

func (h *WaitHandler) Execute(ctx *models.Context, step *models.Step) (result *models.ExecutionResult, err error) {
	duration := step.Config["duration_ms"].(int)
	wakeUp := time.Now().Add(time.Duration(duration) * time.Millisecond)
	return &models.ExecutionResult{
		Status:   models.WAITING,
		NextStep: step.Config["next"].(string),
		WakeUpAt: &wakeUp,
	}, nil
}

func (h *WaitHandler) Validate(step *models.Step) error {
	return nil
}
