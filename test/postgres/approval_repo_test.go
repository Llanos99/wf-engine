package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store/postgres"
	"github.com/pashagolub/pgxmock/v4"
)

func testApproval(id, instanceID, nodeID string) *domain.ApprovalRequest {
	return &domain.ApprovalRequest{
		ID:          id,
		InstanceID:  instanceID,
		NodeID:      nodeID,
		Title:       "Test Approval",
		Description: "Test approval description",
		Approvers:   []string{"manager@test.com"},
		Status:      domain.ApprovalPending,
		CreatedAt:   time.Now(),
	}
}

func TestApprovalRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewApprovalRepo(mock)
	approval := testApproval("approval-001", "instance-001", "approval_node")

	mock.ExpectExec("INSERT INTO approval_requests").
		WithArgs(
			approval.ID,
			approval.InstanceID,
			approval.NodeID,
			approval.Title,
			approval.Description,
			approval.Approvers,
			string(approval.Status),
			pgxmock.AnyArg(), // created_at
			approval.ExpiresAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	ctx := context.Background()
	err = repo.Create(ctx, approval)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestApprovalRepo_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewApprovalRepo(mock)
	now := time.Now()

	mock.ExpectQuery("SELECT id, instance_id").
		WithArgs("approval-001").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "instance_id", "node_id", "title", "description", "approvers",
			"status", "decided_by", "decision_notes", "created_at", "expires_at", "decided_at",
		}).AddRow(
			"approval-001", "instance-001", "approval_node", "Test Approval", "Description",
			[]string{"manager@test.com"}, "pending", nil, nil, now, nil, nil,
		))

	ctx := context.Background()
	result, err := repo.Get(ctx, "approval-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.ID != "approval-001" {
		t.Errorf("expected ID approval-001, got %s", result.ID)
	}
	if result.Status != domain.ApprovalPending {
		t.Errorf("expected status pending, got %s", result.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestApprovalRepo_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewApprovalRepo(mock)
	approval := testApproval("approval-001", "instance-001", "approval_node")
	approval.Approve("manager@test.com", "Looks good")

	mock.ExpectExec("UPDATE approval_requests SET").
		WithArgs(
			approval.ID,
			string(approval.Status),
			approval.DecidedBy,
			approval.Notes,
			pgxmock.AnyArg(), // decided_at
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	ctx := context.Background()
	err = repo.Update(ctx, approval)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestApprovalRepo_ListPending(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewApprovalRepo(mock)
	now := time.Now()

	mock.ExpectQuery("SELECT id, instance_id").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "instance_id", "node_id", "title", "description", "approvers",
			"status", "decided_by", "decision_notes", "created_at", "expires_at", "decided_at",
		}).
			AddRow("approval-001", "instance-001", "node1", "Title 1", "", []string{}, "pending", nil, nil, now, nil, nil).
			AddRow("approval-002", "instance-002", "node2", "Title 2", "", []string{}, "pending", nil, nil, now, nil, nil))

	ctx := context.Background()
	results, err := repo.ListPending(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	for _, r := range results {
		if r.Status != domain.ApprovalPending {
			t.Errorf("expected status pending, got %s", r.Status)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestApprovalRepo_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := postgres.NewApprovalRepo(mock)

	mock.ExpectExec("DELETE FROM approval_requests").
		WithArgs("approval-001").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ctx := context.Background()
	err = repo.Delete(ctx, "approval-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestApprovalRequest_Approve(t *testing.T) {
	approval := domain.NewApprovalRequest("id", "instance", "node", "Title")

	if approval.Status != domain.ApprovalPending {
		t.Errorf("expected pending status, got %s", approval.Status)
	}

	approval.Approve("manager", "Approved")

	if approval.Status != domain.ApprovalApproved {
		t.Errorf("expected approved status, got %s", approval.Status)
	}
	if approval.DecidedBy != "manager" {
		t.Errorf("expected decided_by manager, got %s", approval.DecidedBy)
	}
	if approval.DecidedAt == nil {
		t.Error("expected decided_at to be set")
	}
}

func TestApprovalRequest_Reject(t *testing.T) {
	approval := domain.NewApprovalRequest("id", "instance", "node", "Title")
	approval.Reject("manager", "Not approved")

	if approval.Status != domain.ApprovalRejected {
		t.Errorf("expected rejected status, got %s", approval.Status)
	}
}

func TestApprovalRequest_IsExpired(t *testing.T) {
	approval := domain.NewApprovalRequest("id", "instance", "node", "Title")

	if approval.IsExpired() {
		t.Error("should not be expired when no expiry set")
	}

	past := time.Now().Add(-time.Hour)
	approval.ExpiresAt = &past

	if !approval.IsExpired() {
		t.Error("should be expired when expires_at is in the past")
	}

	future := time.Now().Add(time.Hour)
	approval.ExpiresAt = &future

	if approval.IsExpired() {
		t.Error("should not be expired when expires_at is in the future")
	}
}
