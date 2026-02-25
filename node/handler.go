package node

import "github.com/Llanos99/wf-engine/internal/domain"

// ExecutionResult es el resultado de ejecutar un nodo
type ExecutionResult struct {
	// NextNodeID es el ID del siguiente nodo a ejecutar
	// Vacío si el workflow debe terminar
	NextNodeID string

	// WaitUntil indica que el workflow debe pausarse hasta esta fecha (ISO8601)
	WaitUntil *string

	// WaitForApproval contiene la solicitud de aprobación pendiente
	WaitForApproval *domain.ApprovalRequest

	// Finished indica que el workflow ha terminado exitosamente
	Finished bool

	// Error indica que hubo un error durante la ejecución
	Error error
}

// Handler define la interface que todo handler de nodo debe implementar
type Handler interface {
	// Execute ejecuta la lógica del nodo
	Execute(ctx *ExecutionContext, node domain.NodeDefinition) (*ExecutionResult, error)
}
