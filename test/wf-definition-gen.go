package test

import (
	"errors"

	"github.com/Llanos99/wf-engine/definition"
	"github.com/Llanos99/wf-engine/internal/domain"
)

func TestWorkflowGenerator(path string, loaderType string) (*domain.WorkflowDefinition, error) {
	loader, err := definition.GetLoader(loaderType)
	if err != nil {
		return nil, errors.New("Unknown loader type")
	}
	return loader.Load(path)

}
