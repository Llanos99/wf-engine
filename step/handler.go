package step

import (
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/runtime"
)

type StepHandler interface {
	Execute(instance *domain.WorkflowInstance, step domain.StepDefinition) (*runtime.ExecutionResult, error)
}
