package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/jackc/pgx/v5"
)

// ApprovalRepo implements store.ApprovalRepository for PostgreSQL
type ApprovalRepo struct {
	pool Pool
}

// NewApprovalRepo creates a new approval repository
func NewApprovalRepo(pool Pool) *ApprovalRepo {
	return &ApprovalRepo{pool: pool}
}

// Create creates a new approval request
func (r *ApprovalRepo) Create(ctx context.Context, approval *domain.ApprovalRequest) error {
	query := `
		INSERT INTO approval_requests
			(id, instance_id, node_id, title, description, approvers, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		approval.ID,
		approval.InstanceID,
		approval.NodeID,
		approval.Title,
		approval.Description,
		approval.Approvers,
		string(approval.Status),
		approval.CreatedAt,
		approval.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}

	return nil
}

// Get retrieves an approval request by ID
func (r *ApprovalRepo) Get(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	query := `
		SELECT id, instance_id, node_id, title, description, approvers, status,
		       decided_by, decision_notes, created_at, expires_at, decided_at
		FROM approval_requests
		WHERE id = $1
	`

	var approval domain.ApprovalRequest
	var status string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&approval.ID,
		&approval.InstanceID,
		&approval.NodeID,
		&approval.Title,
		&approval.Description,
		&approval.Approvers,
		&status,
		&approval.DecidedBy,
		&approval.Notes,
		&approval.CreatedAt,
		&approval.ExpiresAt,
		&approval.DecidedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("approval not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get approval: %w", err)
	}

	approval.Status = domain.ApprovalStatus(status)
	return &approval, nil
}

// GetByInstance retrieves all approval requests for an instance
func (r *ApprovalRepo) GetByInstance(ctx context.Context, instanceID string) ([]*domain.ApprovalRequest, error) {
	query := `
		SELECT id, instance_id, node_id, title, description, approvers, status,
		       decided_by, decision_notes, created_at, expires_at, decided_at
		FROM approval_requests
		WHERE instance_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvals: %w", err)
	}
	defer rows.Close()

	return r.scanApprovals(rows)
}

// GetPendingByInstance retrieves the pending approval for an instance
func (r *ApprovalRepo) GetPendingByInstance(ctx context.Context, instanceID string) (*domain.ApprovalRequest, error) {
	query := `
		SELECT id, instance_id, node_id, title, description, approvers, status,
		       decided_by, decision_notes, created_at, expires_at, decided_at
		FROM approval_requests
		WHERE instance_id = $1 AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var approval domain.ApprovalRequest
	var status string

	err := r.pool.QueryRow(ctx, query, instanceID).Scan(
		&approval.ID,
		&approval.InstanceID,
		&approval.NodeID,
		&approval.Title,
		&approval.Description,
		&approval.Approvers,
		&status,
		&approval.DecidedBy,
		&approval.Notes,
		&approval.CreatedAt,
		&approval.ExpiresAt,
		&approval.DecidedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No pending approval
		}
		return nil, fmt.Errorf("failed to get pending approval: %w", err)
	}

	approval.Status = domain.ApprovalStatus(status)
	return &approval, nil
}

// Update updates an approval request
func (r *ApprovalRepo) Update(ctx context.Context, approval *domain.ApprovalRequest) error {
	query := `
		UPDATE approval_requests SET
			status = $2,
			decided_by = $3,
			decision_notes = $4,
			decided_at = $5
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		approval.ID,
		string(approval.Status),
		approval.DecidedBy,
		approval.Notes,
		approval.DecidedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("approval not found: %s", approval.ID)
	}

	return nil
}

// ListPending lists all pending approval requests
func (r *ApprovalRepo) ListPending(ctx context.Context) ([]*domain.ApprovalRequest, error) {
	query := `
		SELECT id, instance_id, node_id, title, description, approvers, status,
		       decided_by, decision_notes, created_at, expires_at, decided_at
		FROM approval_requests
		WHERE status = 'pending'
		  AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending approvals: %w", err)
	}
	defer rows.Close()

	return r.scanApprovals(rows)
}

// Delete deletes an approval request
func (r *ApprovalRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM approval_requests WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete approval: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("approval not found: %s", id)
	}

	return nil
}

// scanApprovals scans multiple approval rows
func (r *ApprovalRepo) scanApprovals(rows pgx.Rows) ([]*domain.ApprovalRequest, error) {
	var approvals []*domain.ApprovalRequest

	for rows.Next() {
		var approval domain.ApprovalRequest
		var status string

		if err := rows.Scan(
			&approval.ID,
			&approval.InstanceID,
			&approval.NodeID,
			&approval.Title,
			&approval.Description,
			&approval.Approvers,
			&status,
			&approval.DecidedBy,
			&approval.Notes,
			&approval.CreatedAt,
			&approval.ExpiresAt,
			&approval.DecidedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}

		approval.Status = domain.ApprovalStatus(status)
		approvals = append(approvals, &approval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating approvals: %w", err)
	}

	return approvals, nil
}
