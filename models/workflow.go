package models

type Workflow struct {
	ID      string
	Name    string
	Steps   []Step
	StartAt string
}

func (wf *Workflow) FindStep(id string) *Step {
	for i := range wf.Steps {
		if wf.Steps[i].ID == id {
			return &wf.Steps[i]
		}
	}
	return nil
}
