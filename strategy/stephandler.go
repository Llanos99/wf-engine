package strategy

import "github.com/Llanos99/wf-engine/models"

type StepHandler interface {
	Execute(ctx *models.Context, step *models.Step) (nextStepID string, err error)
}

var StepHandlers = map[models.StepType]StepHandler{
	models.StepIf:     &ConditionalHandler{},
	models.StepAction: &ActionHandler{},
}
