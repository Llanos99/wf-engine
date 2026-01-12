package strategy

import "github.com/Llanos99/wf-engine/models"

type StepHandler interface {
	Validate(step *models.Step) error
	Execute(ctx *models.Context, step *models.Step) (executionResult *models.ExecutionResult, err error)
}

var StepHandlers = map[models.StepType]StepHandler{
	models.StepIf:     &ConditionalHandler{},
	models.StepAction: &ActionHandler{},
	models.StepWait:   &WaitHandler{},
}
