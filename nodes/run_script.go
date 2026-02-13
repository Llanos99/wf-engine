package nodes

import (
	"encoding/json"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// RunScriptHandler maneja nodos de tipo "run_script"
type RunScriptHandler struct{}

func NewRunScriptHandler() *RunScriptHandler {
	return &RunScriptHandler{}
}

func (h *RunScriptHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.RunScriptSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid run_script spec: %w", err)
	}

	if spec.Script == "" {
		return nil, fmt.Errorf("run_script node '%s' has no script", n.ID)
	}

	// Ejecutar el script
	ctx.Instance.AddLog("debug", n.ID, "Executing script")
	_, err := ctx.ScriptRunner.Execute(spec.Script, ctx)
	if err != nil {
		ctx.Instance.AddLog("error", n.ID, fmt.Sprintf("Script error: %v", err))
		return nil, fmt.Errorf("script execution failed: %w", err)
	}

	ctx.Instance.AddLog("debug", n.ID, "Script executed successfully")

	return &node.ExecutionResult{
		NextNodeID: spec.Next,
	}, nil
}
