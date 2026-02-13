package domain

import (
	"encoding/json"
	"time"
)

type WorkflowStatus string

const (
	StatusPending  WorkflowStatus = "pending"  // Creado pero no iniciado
	StatusRunning  WorkflowStatus = "running"  // En ejecución
	StatusWaiting  WorkflowStatus = "waiting"  // Esperando (timer, approval, etc)
	StatusFailed   WorkflowStatus = "failed"   // Error
	StatusFinished WorkflowStatus = "finished" // Completado exitosamente
)

// WorkflowInstance representa una ejecución específica de un workflow
type WorkflowInstance struct {
	ID           string         `json:"id"`
	DefinitionID string         `json:"definition_id"`
	Version      int            `json:"version"`
	Status       WorkflowStatus `json:"status"`
	CurrentNode  string         `json:"current_node"`
	Variables    Variables      `json:"variables"`
	Logs         []LogEntry     `json:"logs"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	FinishedAt   *time.Time     `json:"finished_at,omitempty"`
}

// LogEntry representa una entrada de log durante la ejecución
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	NodeID    string    `json:"node_id"`
	Message   string    `json:"message"`
}

// NewWorkflowInstance crea una nueva instancia a partir de una definición
func NewWorkflowInstance(id string, def *WorkflowDefinition) *WorkflowInstance {
	now := time.Now()

	// Copiar variables iniciales
	vars := NewVariables()
	if def.InitialVariables != nil {
		for k, v := range def.InitialVariables {
			vars[k] = v
		}
	}

	// Encontrar el nodo Start
	startNode, ok := def.GetStartNode()
	currentNode := ""
	if ok {
		currentNode = startNode.ID
	}

	return &WorkflowInstance{
		ID:           id,
		DefinitionID: def.ID,
		Version:      def.Version,
		Status:       StatusPending,
		CurrentNode:  currentNode,
		Variables:    vars,
		Logs:         make([]LogEntry, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// SetStatus actualiza el estado del workflow
func (w *WorkflowInstance) SetStatus(status WorkflowStatus) {
	w.Status = status
	w.UpdatedAt = time.Now()
	if status == StatusFinished || status == StatusFailed {
		now := time.Now()
		w.FinishedAt = &now
	}
}

// MoveTo mueve la ejecución a otro nodo
func (w *WorkflowInstance) MoveTo(nodeID string) {
	w.CurrentNode = nodeID
	w.UpdatedAt = time.Now()
}

// AddLog agrega una entrada de log
func (w *WorkflowInstance) AddLog(level, nodeID, message string) {
	w.Logs = append(w.Logs, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		NodeID:    nodeID,
		Message:   message,
	})
}

// ToJSON serializa la instancia a JSON
func (w *WorkflowInstance) ToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "  ")
}
