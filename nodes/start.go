package nodes

import (
	"encoding/json"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// StartHandler maneja nodos de tipo "start"
type StartHandler struct{}

func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

func (h *StartHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.StartSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid start spec: %w", err)
	}

	// El nodo Start simplemente inicia el workflow y va al siguiente nodo
	ctx.Instance.SetStatus(domain.StatusRunning)
	ctx.Instance.AddLog("info", n.ID, "Workflow started")

	return &node.ExecutionResult{
		NextNodeID: spec.Next,
	}, nil
}
