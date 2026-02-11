package step

import (
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type StepHandler interface {
	Validate(step *domain.Step) error
	Execute(ctx *runtime.Context, step *domain.Step) (executionResult *runtime.ExecutionResult, err error)
}

var StepHandlers = map[domain.StepType]StepHandler{
	domain.StepIf:     &ConditionalHandler{},
	domain.StepAction: &ActionHandler{},
	domain.StepWait:   &WaitHandler{},
}
