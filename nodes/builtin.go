package nodes

import (
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// RegisterBuiltin registra todos los handlers de nodos builtin
func RegisterBuiltin(r *node.Registry) {
	r.Register(domain.NodeTypeStart, NewStartHandler())
	r.Register(domain.NodeTypeEnd, NewEndHandler())
	r.Register(domain.NodeTypeRunScript, NewRunScriptHandler())
	r.Register(domain.NodeTypeSetValues, NewSetValuesHandler())
	r.Register(domain.NodeTypeConditional, NewConditionalHandler())
	r.Register(domain.NodeTypeWait, NewWaitHandler())
	r.Register(domain.NodeTypeLog, NewLogHandler())
	r.Register(domain.NodeTypeApproval, NewApprovalHandler())
}

// NewRegistryWithBuiltin crea un registry con todos los handlers builtin
func NewRegistryWithBuiltin() *node.Registry {
	r := node.NewRegistry()
	RegisterBuiltin(r)
	return r
}
