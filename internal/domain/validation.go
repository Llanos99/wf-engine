package domain

import "fmt"

func (w *WorkflowDefinition) Validate() error {
	if w.ID == "" {
		return fmt.Errorf("Workflow ID cannot be empty")
	}
	if w.Version <= 0 {
		return fmt.Errorf("Version must be > 0")
	}
	if _, ok := w.Steps[w.StartAt]; !ok {
		return fmt.Errorf("Start At step not found")
	}
	for id, step := range w.Steps {
		if step.ID != id {
			return fmt.Errorf("Step key and step ID mismatch for %s", id)
		}
	}
	return nil
}
