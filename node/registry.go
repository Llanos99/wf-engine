package node

import (
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
)

// Registry mantiene los handlers para cada tipo de nodo
type Registry struct {
	handlers map[domain.NodeType]Handler
}

// NewRegistry crea un nuevo registry vacío
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[domain.NodeType]Handler),
	}
}

// Register registra un handler para un tipo de nodo
func (r *Registry) Register(nodeType domain.NodeType, handler Handler) {
	r.handlers[nodeType] = handler
}

// Get obtiene el handler para un tipo de nodo
func (r *Registry) Get(nodeType domain.NodeType) (Handler, error) {
	handler, ok := r.handlers[nodeType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for node type: %s", nodeType)
	}
	return handler, nil
}

// Has verifica si existe un handler para un tipo de nodo
func (r *Registry) Has(nodeType domain.NodeType) bool {
	_, ok := r.handlers[nodeType]
	return ok
}
