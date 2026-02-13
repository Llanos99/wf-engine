package nodes

import (
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// EndHandler maneja nodos de tipo "end"
type EndHandler struct{}

func NewEndHandler() *EndHandler {
	return &EndHandler{}
}

func (h *EndHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	// El nodo End termina el workflow exitosamente
	ctx.Instance.AddLog("info", n.ID, "Workflow finished")

	return &node.ExecutionResult{
		Finished: true,
	}, nil
}
