package main

import (
	"github.com/Llanos99/wf-engine/definition"
	"github.com/Llanos99/wf-engine/engine"
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/step"
)

func main() {
	loader, err := definition.GetLoader("yaml")
	if err != nil {
		panic(err)
	}
	def, err := loader.Load("testdata/collatz.yaml")
	if err != nil {
		panic(err)
	}
	instance := domain.NewWorkflowInstance("collatz-1", def)

	registry := step.NewRegistry()
	registry.Register(domain.StepTypeIf, step.NewIfHandler())
	registry.Register(domain.StepTypeAction, step.NewActionHandler())

	executor := engine.NewExecutor(registry)

	err = executor.Execute(def, instance)
	if err != nil {
		panic(err)
	}
}
