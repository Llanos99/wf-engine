package postgres_test

import (
	"encoding/json"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
)

// mustJSON marshals a value to JSON, panicking on error (for tests only)
func mustJSON(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// testDefinition creates a sample workflow definition for testing
func testDefinition(id string, version int) *domain.WorkflowDefinition {
	return &domain.WorkflowDefinition{
		ID:          id,
		Version:     version,
		Name:        "Test Workflow",
		Description: "A test workflow",
		Nodes: map[string]domain.NodeDefinition{
			"start": {
				ID:   "start",
				Type: domain.NodeTypeStart,
				Spec: mustJSON(map[string]interface{}{"next": "end"}),
			},
			"end": {
				ID:   "end",
				Type: domain.NodeTypeEnd,
				Spec: mustJSON(map[string]interface{}{}),
			},
		},
		InitialVariables: domain.Variables{
			"test_var": "test_value",
		},
	}
}

// testInstance creates a sample workflow instance for testing
func testInstance(id string, def *domain.WorkflowDefinition) *domain.WorkflowInstance {
	now := time.Now()
	return &domain.WorkflowInstance{
		ID:           id,
		DefinitionID: def.ID,
		Version:      def.Version,
		Status:       domain.StatusPending,
		CurrentNode:  "start",
		Variables: domain.Variables{
			"test_var": "test_value",
		},
		Logs:      []domain.LogEntry{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// testLogEntry creates a sample log entry for testing
func testLogEntry(nodeID, message string) domain.LogEntry {
	return domain.LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		NodeID:    nodeID,
		Message:   message,
	}
}
