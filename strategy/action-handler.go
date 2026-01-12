package strategy

import "github.com/Llanos99/wf-engine/models"

type ActionHandler struct{}

func (h *ActionHandler) Execute(ctx *models.Context, step *models.Step) (executionResult *models.ExecutionResult, err error) {
	return &models.ExecutionResult{
		Status:   models.COMPLETED,
		NextStep: step.Config["next"].(string),
	}, nil
}

func (h *ActionHandler) Validate(step *models.Step) error {
	return nil
}
