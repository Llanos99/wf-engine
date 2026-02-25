package engine

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

const MAX_ITERATIONS = 1000 // Límite de seguridad para evitar loops infinitos

// ExecutionOutput contains the result of workflow execution
type ExecutionOutput struct {
	// PendingApproval is set when the workflow is waiting for approval
	PendingApproval *domain.ApprovalRequest
}

// Executor ejecuta workflows
type Executor struct {
	registry     *node.Registry
	scriptRunner node.ScriptRunner
}

// NewExecutor crea un nuevo executor
func NewExecutor(registry *node.Registry, scriptRunner node.ScriptRunner) *Executor {
	return &Executor{
		registry:     registry,
		scriptRunner: scriptRunner,
	}
}

// Execute ejecuta un workflow desde el inicio hasta el final
func (e *Executor) Execute(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance) (*ExecutionOutput, error) {
	return e.run(def, instance)
}

// Resume continues execution of a paused workflow
func (e *Executor) Resume(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance) (*ExecutionOutput, error) {
	if instance.Status != domain.StatusWaiting {
		return nil, fmt.Errorf("cannot resume workflow with status '%s', expected 'waiting'", instance.Status)
	}

	// Set status back to running
	instance.SetStatus(domain.StatusRunning)
	instance.AddLog("info", instance.CurrentNode, "Workflow resumed")

	return e.run(def, instance)
}

// ResumeAfterApproval continues execution after an approval decision
func (e *Executor) ResumeAfterApproval(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance, approval *domain.ApprovalRequest) (*ExecutionOutput, error) {
	if instance.Status != domain.StatusWaiting {
		return nil, fmt.Errorf("cannot resume workflow with status '%s', expected 'waiting'", instance.Status)
	}

	// Get the approval node spec to determine next node
	nodeDef, ok := def.GetNode(approval.NodeID)
	if !ok {
		return nil, fmt.Errorf("approval node not found: %s", approval.NodeID)
	}

	var spec domain.ApprovalSpec
	if err := parseSpec(nodeDef.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid approval spec: %w", err)
	}

	// Determine next node based on approval status
	var nextNode string
	switch approval.Status {
	case domain.ApprovalApproved:
		nextNode = spec.ApprovedNext
		instance.AddLog("info", approval.NodeID, fmt.Sprintf("Approved by %s", approval.DecidedBy))
	case domain.ApprovalRejected:
		if spec.RejectedNext != "" {
			nextNode = spec.RejectedNext
			instance.AddLog("info", approval.NodeID, fmt.Sprintf("Rejected by %s", approval.DecidedBy))
		} else {
			// No rejected path, fail the workflow
			instance.SetStatus(domain.StatusFailed)
			instance.AddLog("error", approval.NodeID, fmt.Sprintf("Rejected by %s: %s", approval.DecidedBy, approval.Notes))
			return &ExecutionOutput{}, nil
		}
	default:
		return nil, fmt.Errorf("cannot resume with approval status '%s'", approval.Status)
	}

	// Store approval info in variables
	instance.Variables["_approval_decision"] = string(approval.Status)
	instance.Variables["_approval_decided_by"] = approval.DecidedBy
	if approval.Notes != "" {
		instance.Variables["_approval_notes"] = approval.Notes
	}

	// Move to next node and resume
	instance.SetStatus(domain.StatusRunning)
	instance.MoveTo(nextNode)

	return e.run(def, instance)
}

// run is the internal execution loop used by both Execute and Resume
func (e *Executor) run(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance) (*ExecutionOutput, error) {
	output := &ExecutionOutput{}

	// Verificar que hay un nodo actual
	if instance.CurrentNode == "" {
		return nil, errors.New("workflow has no start node")
	}

	// Crear contexto de ejecución
	ctx := node.NewExecutionContext(instance, def, e.scriptRunner)

	iterations := 0

	// Loop de ejecución
	for instance.Status == domain.StatusPending || instance.Status == domain.StatusRunning {
		iterations++
		if iterations > MAX_ITERATIONS {
			instance.SetStatus(domain.StatusFailed)
			return nil, fmt.Errorf("max iterations (%d) exceeded - possible infinite loop", MAX_ITERATIONS)
		}

		// Obtener el nodo actual
		nodeDef, ok := def.GetNode(instance.CurrentNode)
		if !ok {
			instance.SetStatus(domain.StatusFailed)
			return nil, fmt.Errorf("node not found: %s", instance.CurrentNode)
		}

		// Obtener el handler para este tipo de nodo
		handler, err := e.registry.Get(nodeDef.Type)
		if err != nil {
			instance.SetStatus(domain.StatusFailed)
			return nil, fmt.Errorf("no handler for node type '%s': %w", nodeDef.Type, err)
		}

		// Ejecutar el nodo
		result, err := handler.Execute(ctx, *nodeDef)
		if err != nil {
			instance.SetStatus(domain.StatusFailed)
			instance.AddLog("error", nodeDef.ID, err.Error())
			return nil, fmt.Errorf("node '%s' execution failed: %w", nodeDef.ID, err)
		}

		// Procesar el resultado
		if result.Error != nil {
			instance.SetStatus(domain.StatusFailed)
			return nil, result.Error
		}

		if result.Finished {
			instance.SetStatus(domain.StatusFinished)
			return output, nil
		}

		if result.WaitUntil != nil {
			instance.SetStatus(domain.StatusWaiting)
			return output, nil
		}

		if result.WaitForApproval != nil {
			instance.SetStatus(domain.StatusWaiting)
			output.PendingApproval = result.WaitForApproval
			return output, nil
		}

		if result.NextNodeID == "" {
			// Sin siguiente nodo, el workflow termina
			instance.SetStatus(domain.StatusFinished)
			return output, nil
		}

		// Mover al siguiente nodo
		instance.MoveTo(result.NextNodeID)
	}

	return output, nil
}

// parseSpec unmarshals node spec JSON into a typed struct
func parseSpec(spec []byte, v interface{}) error {
	return json.Unmarshal(spec, v)
}
