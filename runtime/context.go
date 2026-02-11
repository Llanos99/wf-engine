package runtime

import (
	"log"

	"github.com/Llanos99/wf-engine/internal/domain"
)

type Context struct {
	WorkflowID     string
	InstanceID     string
	Variables      *domain.Variables
	StepResults    map[string]interface{}
	Logger         *log.Logger
	ExecutionCount map[string]int
}
