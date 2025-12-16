package models

import (
	"errors"
	"fmt"
)

type Workflow struct {
	ID      string
	Name    string
	Steps   map[string]*Step // for branching & loops
	StartAt string
}

func (wf *Workflow) FindStepByID(id string) *Step {
	if step, ok := wf.Steps[id]; ok {
		return step
	}
	return nil
}

func (wf *Workflow) Validate() error {
	if wf.StartAt == "" {
		return errors.New("No starting point was found")
	}
	if len(wf.Steps) == 0 {
		return errors.New("No steps were found")
	}
	startStep, ok := wf.Steps[wf.StartAt]
	if !ok {
		return fmt.Errorf("Start step '%s' does not exists", wf.StartAt)
	}
	if startStep == nil {
		return errors.New("Start step is nil")
	}
	// TODO: How to add loops but avoide the infinite ones?
	return nil
}
