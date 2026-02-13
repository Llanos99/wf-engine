package engine

import (
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

const MAX_ITERATIONS = 1000 // Límite de seguridad para evitar loops infinitos

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
func (e *Executor) Execute(def *domain.WorkflowDefinition, instance *domain.WorkflowInstance) error {
	// Verificar que hay un nodo actual
	if instance.CurrentNode == "" {
		return errors.New("workflow has no start node")
	}

	// Crear contexto de ejecución
	ctx := node.NewExecutionContext(instance, def, e.scriptRunner)

	iterations := 0

	// Loop de ejecución
	for instance.Status == domain.StatusPending || instance.Status == domain.StatusRunning {
		iterations++
		if iterations > MAX_ITERATIONS {
			instance.SetStatus(domain.StatusFailed)
			return fmt.Errorf("max iterations (%d) exceeded - possible infinite loop", MAX_ITERATIONS)
		}

		// Obtener el nodo actual
		nodeDef, ok := def.GetNode(instance.CurrentNode)
		if !ok {
			instance.SetStatus(domain.StatusFailed)
			return fmt.Errorf("node not found: %s", instance.CurrentNode)
		}

		// Obtener el handler para este tipo de nodo
		handler, err := e.registry.Get(nodeDef.Type)
		if err != nil {
			instance.SetStatus(domain.StatusFailed)
			return fmt.Errorf("no handler for node type '%s': %w", nodeDef.Type, err)
		}

		// Ejecutar el nodo
		result, err := handler.Execute(ctx, *nodeDef)
		if err != nil {
			instance.SetStatus(domain.StatusFailed)
			instance.AddLog("error", nodeDef.ID, err.Error())
			return fmt.Errorf("node '%s' execution failed: %w", nodeDef.ID, err)
		}

		// Procesar el resultado
		if result.Error != nil {
			instance.SetStatus(domain.StatusFailed)
			return result.Error
		}

		if result.Finished {
			instance.SetStatus(domain.StatusFinished)
			return nil
		}

		if result.WaitUntil != nil {
			instance.SetStatus(domain.StatusWaiting)
			// En un sistema real, aquí guardaríamos el estado y saldríamos
			// Por ahora continuamos
			return nil
		}

		if result.NextNodeID == "" {
			// Sin siguiente nodo, el workflow termina
			instance.SetStatus(domain.StatusFinished)
			return nil
		}

		// Mover al siguiente nodo
		instance.MoveTo(result.NextNodeID)
	}

	return nil
}
