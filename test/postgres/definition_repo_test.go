package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Llanos99/wf-engine/store/postgres"
	"github.com/pashagolub/pgxmock/v4"
)

func TestDefinitionRepo_Save(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)
	def := testDefinition("test-workflow", 1)

	defJSON, _ := json.Marshal(def)

	mock.ExpectExec("INSERT INTO workflow_definitions").
		WithArgs(def.ID, def.Version, def.Name, def.Description, defJSON).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	ctx := context.Background()
	err = repo.Save(ctx, def)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDefinitionRepo_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)
	def := testDefinition("test-workflow", 1)
	defJSON, _ := json.Marshal(def)

	mock.ExpectQuery("SELECT definition FROM workflow_definitions").
		WithArgs(def.ID, def.Version).
		WillReturnRows(pgxmock.NewRows([]string{"definition"}).AddRow(defJSON))

	ctx := context.Background()
	result, err := repo.Get(ctx, def.ID, def.Version)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.ID != def.ID {
		t.Errorf("expected ID %s, got %s", def.ID, result.ID)
	}
	if result.Version != def.Version {
		t.Errorf("expected version %d, got %d", def.Version, result.Version)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDefinitionRepo_GetLatest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)
	def := testDefinition("test-workflow", 2)
	defJSON, _ := json.Marshal(def)

	mock.ExpectQuery("SELECT definition FROM workflow_definitions").
		WithArgs(def.ID).
		WillReturnRows(pgxmock.NewRows([]string{"definition"}).AddRow(defJSON))

	ctx := context.Background()
	result, err := repo.GetLatest(ctx, def.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.Version != 2 {
		t.Errorf("expected version 2, got %d", result.Version)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDefinitionRepo_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)

	def1 := testDefinition("workflow-1", 1)
	def2 := testDefinition("workflow-2", 1)
	def1JSON, _ := json.Marshal(def1)
	def2JSON, _ := json.Marshal(def2)

	mock.ExpectQuery("SELECT DISTINCT ON").
		WillReturnRows(pgxmock.NewRows([]string{"definition"}).
			AddRow(def1JSON).
			AddRow(def2JSON))

	ctx := context.Background()
	results, err := repo.List(ctx)
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

func TestDefinitionRepo_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)

	mock.ExpectExec("DELETE FROM workflow_definitions").
		WithArgs("test-workflow", 1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ctx := context.Background()
	err = repo.Delete(ctx, "test-workflow", 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDefinitionRepo_Delete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewDefinitionRepo(mock)

	mock.ExpectExec("DELETE FROM workflow_definitions").
		WithArgs("nonexistent", 1).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	ctx := context.Background()
	err = repo.Delete(ctx, "nonexistent", 1)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
