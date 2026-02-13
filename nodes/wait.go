package nodes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// WaitHandler maneja nodos de tipo "wait"
type WaitHandler struct{}

func NewWaitHandler() *WaitHandler {
	return &WaitHandler{}
}

func (h *WaitHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.WaitSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid wait spec: %w", err)
	}

	duration, err := time.ParseDuration(spec.Duration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration '%s': %w", spec.Duration, err)
	}

	// Calcular el tiempo de espera
	waitUntil := time.Now().Add(duration).Format(time.RFC3339)

	ctx.Instance.AddLog("info", n.ID, fmt.Sprintf("Waiting for %s (until %s)", spec.Duration, waitUntil))

	// Por ahora, para simplificar, esperamos sincrónicamente
	// En un sistema real, esto guardaría el estado y retomaría después
	time.Sleep(duration)

	ctx.Instance.AddLog("info", n.ID, "Wait completed")

	return &node.ExecutionResult{
		NextNodeID: spec.Next,
	}, nil
}
