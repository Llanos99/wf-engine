package node

import (
	"github.com/Llanos99/wf-engine/internal/domain"
)

// ExecutionContext contiene todo lo necesario para ejecutar un nodo
type ExecutionContext struct {
	// Instance es la instancia del workflow en ejecución
	Instance *domain.WorkflowInstance

	// Definition es la definición del workflow
	Definition *domain.WorkflowDefinition

	// ScriptRunner es el ejecutor de scripts JavaScript (se inyecta)
	ScriptRunner ScriptRunner
}

// ScriptRunner es la interface para ejecutar scripts
type ScriptRunner interface {
	// Execute ejecuta un script JavaScript y retorna el resultado
	Execute(script string, ctx *ExecutionContext) (interface{}, error)

	// ExecuteBool ejecuta un script y espera un boolean
	ExecuteBool(script string, ctx *ExecutionContext) (bool, error)
}

// NewExecutionContext crea un nuevo contexto de ejecución
func NewExecutionContext(instance *domain.WorkflowInstance, definition *domain.WorkflowDefinition, runner ScriptRunner) *ExecutionContext {
	return &ExecutionContext{
		Instance:     instance,
		Definition:   definition,
		ScriptRunner: runner,
	}
}
