package nodes

import (
	"encoding/json"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// SetValuesHandler maneja nodos de tipo "set_values"
type SetValuesHandler struct{}

func NewSetValuesHandler() *SetValuesHandler {
	return &SetValuesHandler{}
}

func (h *SetValuesHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.SetValuesSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid set_values spec: %w", err)
	}

	// Asignar cada valor a las variables
	for key, value := range spec.Assignments {
		ctx.Instance.Variables[key] = value
		ctx.Instance.AddLog("debug", n.ID, fmt.Sprintf("Set %s = %v", key, value))
	}

	return &node.ExecutionResult{
		NextNodeID: spec.Next,
	}, nil
}
