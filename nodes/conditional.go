package nodes

import (
	"encoding/json"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// ConditionalHandler maneja nodos de tipo "conditional"
type ConditionalHandler struct{}

func NewConditionalHandler() *ConditionalHandler {
	return &ConditionalHandler{}
}

func (h *ConditionalHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.ConditionalSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid conditional spec: %w", err)
	}

	var result bool
	var err error

	// Evaluar usando Script o Expression
	if spec.Script != "" {
		// Usar script JavaScript para evaluar
		result, err = ctx.ScriptRunner.ExecuteBool(spec.Script, ctx)
		if err != nil {
			return nil, fmt.Errorf("conditional script error: %w", err)
		}
	} else if spec.Expression != nil {
		// Usar expresión simple
		result, err = spec.Expression.Evaluate(ctx.Instance.Variables)
		if err != nil {
			return nil, fmt.Errorf("conditional expression error: %w", err)
		}
	} else {
		return nil, fmt.Errorf("conditional node '%s' has no script or expression", n.ID)
	}

	// Determinar el siguiente nodo basado en el resultado
	var nextNode string
	if result {
		nextNode = spec.TrueNext
		ctx.Instance.AddLog("debug", n.ID, "Condition evaluated to true")
	} else {
		nextNode = spec.FalseNext
		ctx.Instance.AddLog("debug", n.ID, "Condition evaluated to false")
	}

	return &node.ExecutionResult{
		NextNodeID: nextNode,
	}, nil
}
