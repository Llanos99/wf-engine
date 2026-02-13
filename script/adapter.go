package script

import (
	"github.com/Llanos99/wf-engine/node"
)

// Adapter implementa node.ScriptRunner usando el Runtime de goja
type Adapter struct {
	runtime *Runtime
}

// NewAdapter crea un nuevo adaptador
func NewAdapter() *Adapter {
	return &Adapter{
		runtime: NewRuntime(),
	}
}

// Execute implementa node.ScriptRunner
func (a *Adapter) Execute(script string, ctx *node.ExecutionContext) (interface{}, error) {
	return a.runtime.Execute(script, &ctx.Instance.Variables, ctx.Instance)
}

// ExecuteBool implementa node.ScriptRunner
func (a *Adapter) ExecuteBool(script string, ctx *node.ExecutionContext) (bool, error) {
	return a.runtime.ExecuteBool(script, &ctx.Instance.Variables, ctx.Instance)
}
