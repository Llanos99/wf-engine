package domain

import "time"

type WorkflowStatus string

const (
	StatusRunning  WorkflowStatus = "running"
	StatusWaiting  WorkflowStatus = "waiting"
	StatusFailed   WorkflowStatus = "failed"
	StatusFinished WorkflowStatus = "finished"
)

type WorkflowInstance struct {
	ID           string         `yaml:"id" json:"id"`
	DefinitionID string         `yaml:"definition_id" json:"definition_id"`
	Version      int            `yaml:"version" json:"version"`
	Status       WorkflowStatus `yaml:"status" json:"status"`
	CurrentStep  string         `yaml:"current_step" json:"current_step"`
	Variables    Variables      `yaml:"variables" json:"variables"`
	CreatedAt    time.Time      `yaml:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `yaml:"updated_at" json:"updated_at"`
	FinishedAt   time.Time      `yaml:"finished_at" json:"finished_at"`
}

func NewWorkflowInstance(id string, def *WorkflowDefinition) *WorkflowInstance {
	now := time.Now()
	vars := NewVariables()
	if def.InitialVariables != nil {
		for k, v := range def.InitialVariables {
			vars[k] = v
		}
	}
	return &WorkflowInstance{
		ID:           id,
		DefinitionID: def.ID,
		Version:      def.Version,
		Status:       StatusRunning,
		CurrentStep:  def.StartAt,
		Variables:    vars,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (w *WorkflowInstance) SetStatus(status WorkflowStatus) {
	w.Status = status
	w.UpdatedAt = time.Now()
	if status == StatusFinished || status == StatusFailed {
		now := time.Now()
		w.FinishedAt = now
	}
}

func (w *WorkflowInstance) MoveTo(stepID string) {
	w.CurrentStep = stepID
	w.UpdatedAt = time.Now()
}
