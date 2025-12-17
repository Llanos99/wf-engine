package strategy

import "github.com/Llanos99/wf-engine/models"

type ActionHandler struct{}

func (h *ActionHandler) Execute(ctx *models.Context, step *models.Step) (nextStepID string, err error) {
	return "", nil
}
