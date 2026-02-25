//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store"
	"github.com/Llanos99/wf-engine/store/postgres"
)

// getTestConfig returns the configuration for integration tests.
// Uses environment variables or defaults to local development settings.
func getTestConfig() postgres.Config {
	cfg := postgres.DefaultConfig()

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		cfg.Host = host
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.User = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.Password = password
	}
	if database := os.Getenv("POSTGRES_DB"); database != "" {
		cfg.Database = database
	}

	return cfg
}

// setupTestStore creates a store and cleans up test data
func setupTestStore(t *testing.T) *postgres.PostgresStore {
	t.Helper()

	ctx := context.Background()
	s, err := postgres.NewStore(ctx, getTestConfig())
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Clean up test data
	_, _ = s.Pool().Exec(ctx, "DELETE FROM execution_logs WHERE instance_id LIKE 'test-%'")
	_, _ = s.Pool().Exec(ctx, "DELETE FROM workflow_instances WHERE id LIKE 'test-%'")
	_, _ = s.Pool().Exec(ctx, "DELETE FROM workflow_definitions WHERE id LIKE 'test-%'")

	return s
}

func TestIntegration_DefinitionRepository(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	ctx := context.Background()
	repo := s.Definitions()

	t.Run("Save and Get", func(t *testing.T) {
		def := testDefinition("test-workflow-1", 1)

		err := repo.Save(ctx, def)
		if err != nil {
			t.Fatalf("failed to save definition: %v", err)
		}

		retrieved, err := repo.Get(ctx, def.ID, def.Version)
		if err != nil {
			t.Fatalf("failed to get definition: %v", err)
		}

		if retrieved.ID != def.ID {
			t.Errorf("expected ID %s, got %s", def.ID, retrieved.ID)
		}
		if retrieved.Name != def.Name {
			t.Errorf("expected Name %s, got %s", def.Name, retrieved.Name)
		}
		if len(retrieved.Nodes) != len(def.Nodes) {
			t.Errorf("expected %d nodes, got %d", len(def.Nodes), len(retrieved.Nodes))
		}
	})

	t.Run("Save Multiple Versions", func(t *testing.T) {
		defV1 := testDefinition("test-workflow-2", 1)
		defV2 := testDefinition("test-workflow-2", 2)
		defV2.Name = "Test Workflow v2"

		err := repo.Save(ctx, defV1)
		if err != nil {
			t.Fatalf("failed to save v1: %v", err)
		}

		err = repo.Save(ctx, defV2)
		if err != nil {
			t.Fatalf("failed to save v2: %v", err)
		}

		latest, err := repo.GetLatest(ctx, "test-workflow-2")
		if err != nil {
			t.Fatalf("failed to get latest: %v", err)
		}

		if latest.Version != 2 {
			t.Errorf("expected version 2, got %d", latest.Version)
		}
		if latest.Name != "Test Workflow v2" {
			t.Errorf("expected name 'Test Workflow v2', got %s", latest.Name)
		}
	})

	t.Run("List", func(t *testing.T) {
		definitions, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("failed to list: %v", err)
		}

		if len(definitions) < 2 {
			t.Errorf("expected at least 2 definitions, got %d", len(definitions))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		def := testDefinition("test-workflow-delete", 1)

		err := repo.Save(ctx, def)
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		err = repo.Delete(ctx, def.ID, def.Version)
		if err != nil {
			t.Fatalf("failed to delete: %v", err)
		}

		_, err = repo.Get(ctx, def.ID, def.Version)
		if err == nil {
			t.Error("expected error after delete, got nil")
		}
	})
}

func TestIntegration_InstanceRepository(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	ctx := context.Background()

	// First, create a definition
	def := testDefinition("test-workflow-inst", 1)
	err := s.Definitions().Save(ctx, def)
	if err != nil {
		t.Fatalf("failed to save definition: %v", err)
	}

	repo := s.Instances()

	t.Run("Create and Get", func(t *testing.T) {
		instance := testInstance("test-instance-1", def)

		err := repo.Create(ctx, instance)
		if err != nil {
			t.Fatalf("failed to create instance: %v", err)
		}

		retrieved, err := repo.Get(ctx, instance.ID)
		if err != nil {
			t.Fatalf("failed to get instance: %v", err)
		}

		if retrieved.ID != instance.ID {
			t.Errorf("expected ID %s, got %s", instance.ID, retrieved.ID)
		}
		if retrieved.Status != instance.Status {
			t.Errorf("expected status %s, got %s", instance.Status, retrieved.Status)
		}
		if retrieved.CurrentNode != instance.CurrentNode {
			t.Errorf("expected current_node %s, got %s", instance.CurrentNode, retrieved.CurrentNode)
		}
	})

	t.Run("Update", func(t *testing.T) {
		instance := testInstance("test-instance-2", def)
		err := repo.Create(ctx, instance)
		if err != nil {
			t.Fatalf("failed to create: %v", err)
		}

		// Update status and node
		instance.SetStatus(domain.StatusRunning)
		instance.MoveTo("process_node")
		instance.Variables["new_var"] = "new_value"

		err = repo.Update(ctx, instance)
		if err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		retrieved, err := repo.Get(ctx, instance.ID)
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}

		if retrieved.Status != domain.StatusRunning {
			t.Errorf("expected status running, got %s", retrieved.Status)
		}
		if retrieved.CurrentNode != "process_node" {
			t.Errorf("expected current_node process_node, got %s", retrieved.CurrentNode)
		}
		if retrieved.Variables["new_var"] != "new_value" {
			t.Error("expected new_var to be set")
		}
	})

	t.Run("List with Filter", func(t *testing.T) {
		// Create instances with different statuses
		running := testInstance("test-instance-running", def)
		running.Status = domain.StatusRunning
		repo.Create(ctx, running)

		waiting := testInstance("test-instance-waiting", def)
		waiting.Status = domain.StatusWaiting
		repo.Create(ctx, waiting)

		finished := testInstance("test-instance-finished", def)
		finished.Status = domain.StatusFinished
		repo.Create(ctx, finished)

		// List running only
		filter := store.InstanceFilter{
			Status: domain.StatusRunning,
		}
		results, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("failed to list: %v", err)
		}

		for _, r := range results {
			if r.Status != domain.StatusRunning {
				t.Errorf("expected only running instances, got %s", r.Status)
			}
		}
	})

	t.Run("SaveLog and GetLogs", func(t *testing.T) {
		instance := testInstance("test-instance-logs", def)
		err := repo.Create(ctx, instance)
		if err != nil {
			t.Fatalf("failed to create: %v", err)
		}

		// Add logs
		logs := []domain.LogEntry{
			{Timestamp: time.Now(), Level: "info", NodeID: "start", Message: "Started"},
			{Timestamp: time.Now().Add(time.Second), Level: "debug", NodeID: "process", Message: "Processing"},
			{Timestamp: time.Now().Add(2 * time.Second), Level: "info", NodeID: "end", Message: "Finished"},
		}

		for _, log := range logs {
			err := repo.SaveLog(ctx, instance.ID, log)
			if err != nil {
				t.Fatalf("failed to save log: %v", err)
			}
		}

		retrieved, err := repo.GetLogs(ctx, instance.ID)
		if err != nil {
			t.Fatalf("failed to get logs: %v", err)
		}

		if len(retrieved) != 3 {
			t.Errorf("expected 3 logs, got %d", len(retrieved))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		instance := testInstance("test-instance-delete", def)
		err := repo.Create(ctx, instance)
		if err != nil {
			t.Fatalf("failed to create: %v", err)
		}

		err = repo.Delete(ctx, instance.ID)
		if err != nil {
			t.Fatalf("failed to delete: %v", err)
		}

		_, err = repo.Get(ctx, instance.ID)
		if err == nil {
			t.Error("expected error after delete, got nil")
		}
	})
}

func TestIntegration_WorkflowPauseResume(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	ctx := context.Background()

	// Create definition with a wait node
	def := &domain.WorkflowDefinition{
		ID:          "test-pause-resume",
		Version:     1,
		Name:        "Pause Resume Test",
		Description: "Test workflow for pause/resume",
		Nodes: map[string]domain.NodeDefinition{
			"start": {ID: "start", Type: domain.NodeTypeStart, Spec: mustJSON(map[string]interface{}{"next": "wait"})},
			"wait":  {ID: "wait", Type: domain.NodeTypeWait, Spec: mustJSON(map[string]interface{}{"duration": "1s", "next": "end"})},
			"end":   {ID: "end", Type: domain.NodeTypeEnd, Spec: mustJSON(map[string]interface{}{})},
		},
		InitialVariables: domain.Variables{"counter": 0},
	}

	err := s.Definitions().Save(ctx, def)
	if err != nil {
		t.Fatalf("failed to save definition: %v", err)
	}

	// Create and start instance
	instance := domain.NewWorkflowInstance("test-pause-resume-inst", def)
	err = s.Instances().Create(ctx, instance)
	if err != nil {
		t.Fatalf("failed to create instance: %v", err)
	}

	// Simulate execution reaching wait node
	instance.SetStatus(domain.StatusRunning)
	instance.MoveTo("wait")
	instance.Variables["counter"] = 1
	instance.AddLog("info", "start", "Workflow started")

	// Pause at wait node
	instance.SetStatus(domain.StatusWaiting)
	instance.AddLog("info", "wait", "Waiting for timer")

	err = s.Instances().Update(ctx, instance)
	if err != nil {
		t.Fatalf("failed to update instance: %v", err)
	}

	// Retrieve and verify paused state
	retrieved, err := s.Instances().Get(ctx, instance.ID)
	if err != nil {
		t.Fatalf("failed to get instance: %v", err)
	}

	if retrieved.Status != domain.StatusWaiting {
		t.Errorf("expected status waiting, got %s", retrieved.Status)
	}
	if retrieved.CurrentNode != "wait" {
		t.Errorf("expected current_node wait, got %s", retrieved.CurrentNode)
	}
	if retrieved.Variables["counter"] != float64(1) { // JSON unmarshals numbers as float64
		t.Errorf("expected counter 1, got %v", retrieved.Variables["counter"])
	}

	// Simulate resume
	retrieved.SetStatus(domain.StatusRunning)
	retrieved.MoveTo("end")
	retrieved.AddLog("info", "wait", "Timer completed, resuming")

	// Complete workflow
	retrieved.SetStatus(domain.StatusFinished)
	retrieved.AddLog("info", "end", "Workflow finished")

	err = s.Instances().Update(ctx, retrieved)
	if err != nil {
		t.Fatalf("failed to update after resume: %v", err)
	}

	// Verify final state
	final, err := s.Instances().Get(ctx, instance.ID)
	if err != nil {
		t.Fatalf("failed to get final instance: %v", err)
	}

	if final.Status != domain.StatusFinished {
		t.Errorf("expected status finished, got %s", final.Status)
	}
	if final.FinishedAt == nil {
		t.Error("expected finished_at to be set")
	}
	if len(final.Logs) != 4 {
		t.Errorf("expected 4 logs, got %d", len(final.Logs))
	}
}

func TestIntegration_Store(t *testing.T) {
	ctx := context.Background()

	t.Run("Connection", func(t *testing.T) {
		s, err := postgres.NewStore(ctx, getTestConfig())
		if err != nil {
			t.Fatalf("failed to connect: %v", err)
		}
		defer s.Close()

		// Verify pool is working
		var result int
		err = s.Pool().QueryRow(ctx, "SELECT 1").Scan(&result)
		if err != nil {
			t.Errorf("failed to query: %v", err)
		}
		if result != 1 {
			t.Errorf("expected 1, got %d", result)
		}
	})

	t.Run("Invalid Connection", func(t *testing.T) {
		cfg := postgres.Config{
			Host:     "invalid-host",
			Port:     5432,
			User:     "invalid",
			Password: "invalid",
			Database: "invalid",
			SSLMode:  "disable",
		}

		_, err := postgres.NewStore(ctx, cfg)
		if err == nil {
			t.Error("expected error for invalid connection, got nil")
		}
	})
}
