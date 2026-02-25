package postgres_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store"
	"github.com/Llanos99/wf-engine/store/postgres"
	"github.com/pashagolub/pgxmock/v4"
)

func TestInstanceRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	def := testDefinition("test-workflow", 1)
	instance := testInstance("instance-001", def)

	varsJSON, _ := json.Marshal(instance.Variables)

	mock.ExpectExec("INSERT INTO workflow_instances").
		WithArgs(
			instance.ID,
			instance.DefinitionID,
			instance.Version,
			string(instance.Status),
			instance.CurrentNode,
			varsJSON,
			pgxmock.AnyArg(), // created_at
			pgxmock.AnyArg(), // updated_at
			instance.FinishedAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	ctx := context.Background()
	err = repo.Create(ctx, instance)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	def := testDefinition("test-workflow", 1)
	instance := testInstance("instance-001", def)
	instance.Status = domain.StatusRunning
	instance.CurrentNode = "process"

	varsJSON, _ := json.Marshal(instance.Variables)

	mock.ExpectExec("UPDATE workflow_instances SET").
		WithArgs(
			instance.ID,
			string(instance.Status),
			instance.CurrentNode,
			varsJSON,
			pgxmock.AnyArg(), // updated_at
			instance.FinishedAt,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	ctx := context.Background()
	err = repo.Update(ctx, instance)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_Update_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	def := testDefinition("test-workflow", 1)
	instance := testInstance("nonexistent", def)

	varsJSON, _ := json.Marshal(instance.Variables)

	mock.ExpectExec("UPDATE workflow_instances SET").
		WithArgs(
			instance.ID,
			string(instance.Status),
			instance.CurrentNode,
			varsJSON,
			pgxmock.AnyArg(),
			instance.FinishedAt,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	ctx := context.Background()
	err = repo.Update(ctx, instance)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	now := time.Now()
	varsJSON, _ := json.Marshal(domain.Variables{"key": "value"})

	// Mock the main query
	mock.ExpectQuery("SELECT id, definition_id").
		WithArgs("instance-001").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "definition_id", "definition_version", "status",
			"current_node", "variables", "created_at", "updated_at", "finished_at",
		}).AddRow(
			"instance-001", "test-workflow", 1, "running",
			"process", varsJSON, now, now, nil,
		))

	// Mock the logs query
	mock.ExpectQuery("SELECT timestamp, level, node_id, message").
		WithArgs("instance-001").
		WillReturnRows(pgxmock.NewRows([]string{"timestamp", "level", "node_id", "message"}).
			AddRow(now, "info", "start", "Workflow started"))

	ctx := context.Background()
	result, err := repo.Get(ctx, "instance-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.ID != "instance-001" {
		t.Errorf("expected ID instance-001, got %s", result.ID)
	}
	if result.Status != domain.StatusRunning {
		t.Errorf("expected status running, got %s", result.Status)
	}
	if len(result.Logs) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(result.Logs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	now := time.Now()
	varsJSON, _ := json.Marshal(domain.Variables{})

	mock.ExpectQuery("SELECT id, definition_id").
		WithArgs("test-workflow", "running", 100, 0).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "definition_id", "definition_version", "status",
			"current_node", "variables", "created_at", "updated_at", "finished_at",
		}).
			AddRow("instance-001", "test-workflow", 1, "running", "node1", varsJSON, now, now, nil).
			AddRow("instance-002", "test-workflow", 1, "running", "node2", varsJSON, now, now, nil))

	ctx := context.Background()
	filter := store.InstanceFilter{
		DefinitionID: "test-workflow",
		Status:       domain.StatusRunning,
	}
	results, err := repo.List(ctx, filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)

	mock.ExpectExec("DELETE FROM workflow_instances").
		WithArgs("instance-001").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ctx := context.Background()
	err = repo.Delete(ctx, "instance-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_SaveLog(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	logEntry := testLogEntry("process", "Processing data")

	mock.ExpectExec("INSERT INTO execution_logs").
		WithArgs("instance-001", pgxmock.AnyArg(), logEntry.Level, logEntry.NodeID, logEntry.Message).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	ctx := context.Background()
	err = repo.SaveLog(ctx, "instance-001", logEntry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_GetLogs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	now := time.Now()

	mock.ExpectQuery("SELECT timestamp, level, node_id, message").
		WithArgs("instance-001").
		WillReturnRows(pgxmock.NewRows([]string{"timestamp", "level", "node_id", "message"}).
			AddRow(now, "info", "start", "Started").
			AddRow(now.Add(time.Second), "info", "process", "Processing").
			AddRow(now.Add(2*time.Second), "info", "end", "Finished"))

	ctx := context.Background()
	logs, err := repo.GetLogs(ctx, "instance-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(logs) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(logs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestInstanceRepo_GetWaitingInstances(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewInstanceRepo(mock)
	now := time.Now()
	varsJSON, _ := json.Marshal(domain.Variables{})

	mock.ExpectQuery("SELECT id, definition_id").
		WithArgs("", "waiting", 100, 0).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "definition_id", "definition_version", "status",
			"current_node", "variables", "created_at", "updated_at", "finished_at",
		}).
			AddRow("instance-001", "workflow-1", 1, "waiting", "wait_node", varsJSON, now, now, nil))

	ctx := context.Background()
	results, err := repo.GetWaitingInstances(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != domain.StatusWaiting {
		t.Errorf("expected status waiting, got %s", results[0].Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
