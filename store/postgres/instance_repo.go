package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store"
	"github.com/jackc/pgx/v5"
)

// InstanceRepo implements store.InstanceRepository for PostgreSQL
type InstanceRepo struct {
	pool Pool
}

// NewInstanceRepo creates a new instance repository
func NewInstanceRepo(pool Pool) *InstanceRepo {
	return &InstanceRepo{pool: pool}
}

// Create creates a new workflow instance
func (r *InstanceRepo) Create(ctx context.Context, instance *domain.WorkflowInstance) error {
	varsJSON, err := json.Marshal(instance.Variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `
		INSERT INTO workflow_instances
			(id, definition_id, definition_version, status, current_node, variables, created_at, updated_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.pool.Exec(ctx, query,
		instance.ID,
		instance.DefinitionID,
		instance.Version,
		string(instance.Status),
		instance.CurrentNode,
		varsJSON,
		instance.CreatedAt,
		instance.UpdatedAt,
		instance.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	// Save any existing logs
	for _, log := range instance.Logs {
		if err := r.SaveLog(ctx, instance.ID, log); err != nil {
			return fmt.Errorf("failed to save log: %w", err)
		}
	}

	return nil
}

// Update updates an existing workflow instance
func (r *InstanceRepo) Update(ctx context.Context, instance *domain.WorkflowInstance) error {
	varsJSON, err := json.Marshal(instance.Variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `
		UPDATE workflow_instances SET
			status = $2,
			current_node = $3,
			variables = $4,
			updated_at = $5,
			finished_at = $6
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		instance.ID,
		string(instance.Status),
		instance.CurrentNode,
		varsJSON,
		instance.UpdatedAt,
		instance.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("instance not found: %s", instance.ID)
	}

	return nil
}

// Get retrieves a workflow instance by ID
func (r *InstanceRepo) Get(ctx context.Context, id string) (*domain.WorkflowInstance, error) {
	query := `
		SELECT id, definition_id, definition_version, status, current_node, variables,
		       created_at, updated_at, finished_at
		FROM workflow_instances
		WHERE id = $1
	`

	var instance domain.WorkflowInstance
	var varsJSON []byte
	var status string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&instance.ID,
		&instance.DefinitionID,
		&instance.Version,
		&status,
		&instance.CurrentNode,
		&varsJSON,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("instance not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	instance.Status = domain.WorkflowStatus(status)

	if err := json.Unmarshal(varsJSON, &instance.Variables); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
	}

	// Load logs
	logs, err := r.GetLogs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	instance.Logs = logs

	return &instance, nil
}

// List retrieves workflow instances with optional filtering
func (r *InstanceRepo) List(ctx context.Context, filter store.InstanceFilter) ([]*domain.WorkflowInstance, error) {
	query := `
		SELECT id, definition_id, definition_version, status, current_node, variables,
		       created_at, updated_at, finished_at
		FROM workflow_instances
		WHERE ($1 = '' OR definition_id = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.pool.Query(ctx, query, filter.DefinitionID, string(filter.Status), limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}
	defer rows.Close()

	var instances []*domain.WorkflowInstance
	for rows.Next() {
		var instance domain.WorkflowInstance
		var varsJSON []byte
		var status string

		if err := rows.Scan(
			&instance.ID,
			&instance.DefinitionID,
			&instance.Version,
			&status,
			&instance.CurrentNode,
			&varsJSON,
			&instance.CreatedAt,
			&instance.UpdatedAt,
			&instance.FinishedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan instance: %w", err)
		}

		instance.Status = domain.WorkflowStatus(status)

		if err := json.Unmarshal(varsJSON, &instance.Variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}

		instances = append(instances, &instance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating instances: %w", err)
	}

	return instances, nil
}

// Delete removes a workflow instance
func (r *InstanceRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM workflow_instances WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("instance not found: %s", id)
	}

	return nil
}

// SaveLog saves a log entry for an instance
func (r *InstanceRepo) SaveLog(ctx context.Context, instanceID string, log domain.LogEntry) error {
	query := `
		INSERT INTO execution_logs (instance_id, timestamp, level, node_id, message)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query, instanceID, log.Timestamp, log.Level, log.NodeID, log.Message)
	if err != nil {
		return fmt.Errorf("failed to save log: %w", err)
	}

	return nil
}

// GetLogs retrieves all logs for an instance
func (r *InstanceRepo) GetLogs(ctx context.Context, instanceID string) ([]domain.LogEntry, error) {
	query := `
		SELECT timestamp, level, node_id, message
		FROM execution_logs
		WHERE instance_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.pool.Query(ctx, query, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.LogEntry
	for rows.Next() {
		var log domain.LogEntry
		if err := rows.Scan(&log.Timestamp, &log.Level, &log.NodeID, &log.Message); err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetWaitingInstances retrieves all instances in waiting status that are ready to resume
func (r *InstanceRepo) GetWaitingInstances(ctx context.Context) ([]*domain.WorkflowInstance, error) {
	return r.List(ctx, store.InstanceFilter{
		Status: domain.StatusWaiting,
	})
}
