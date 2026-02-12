package step

import "github.com/Llanos99/wf-engine/internal/domain"

type Registry struct {
	handlers map[domain.StepType]StepHandler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[domain.StepType]StepHandler),
	}
}

func (r *Registry) Register(stepType domain.StepType, handler StepHandler) {
	r.handlers[stepType] = handler
}

func (r *Registry) Get(stepType domain.StepType) StepHandler {
	return r.handlers[stepType]
}
