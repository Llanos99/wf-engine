package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/jackc/pgx/v5"
)

// DefinitionRepo implements store.DefinitionRepository for PostgreSQL
type DefinitionRepo struct {
	pool Pool
}

// NewDefinitionRepo creates a new definition repository
func NewDefinitionRepo(pool Pool) *DefinitionRepo {
	return &DefinitionRepo{pool: pool}
}

// Save saves a workflow definition to the database
func (r *DefinitionRepo) Save(ctx context.Context, def *domain.WorkflowDefinition) error {
	// Serialize the full definition to JSON
	defJSON, err := json.Marshal(def)
	if err != nil {
		return fmt.Errorf("failed to marshal definition: %w", err)
	}

	query := `
		INSERT INTO workflow_definitions (id, version, name, description, definition)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id, version) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			definition = EXCLUDED.definition
	`

	_, err = r.pool.Exec(ctx, query, def.ID, def.Version, def.Name, def.Description, defJSON)
	if err != nil {
		return fmt.Errorf("failed to save definition: %w", err)
	}

	return nil
}

// Get retrieves a workflow definition by ID and version
func (r *DefinitionRepo) Get(ctx context.Context, id string, version int) (*domain.WorkflowDefinition, error) {
	query := `
		SELECT definition FROM workflow_definitions
		WHERE id = $1 AND version = $2
	`

	var defJSON []byte
	err := r.pool.QueryRow(ctx, query, id, version).Scan(&defJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("definition not found: %s v%d", id, version)
		}
		return nil, fmt.Errorf("failed to get definition: %w", err)
	}

	var def domain.WorkflowDefinition
	if err := json.Unmarshal(defJSON, &def); err != nil {
		return nil, fmt.Errorf("failed to unmarshal definition: %w", err)
	}

	return &def, nil
}

// GetLatest retrieves the latest version of a workflow definition
func (r *DefinitionRepo) GetLatest(ctx context.Context, id string) (*domain.WorkflowDefinition, error) {
	query := `
		SELECT definition FROM workflow_definitions
		WHERE id = $1
		ORDER BY version DESC
		LIMIT 1
	`

	var defJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&defJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("definition not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get latest definition: %w", err)
	}

	var def domain.WorkflowDefinition
	if err := json.Unmarshal(defJSON, &def); err != nil {
		return nil, fmt.Errorf("failed to unmarshal definition: %w", err)
	}

	return &def, nil
}

// List returns all workflow definitions (latest version of each)
func (r *DefinitionRepo) List(ctx context.Context) ([]*domain.WorkflowDefinition, error) {
	query := `
		SELECT DISTINCT ON (id) definition
		FROM workflow_definitions
		ORDER BY id, version DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list definitions: %w", err)
	}
	defer rows.Close()

	var definitions []*domain.WorkflowDefinition
	for rows.Next() {
		var defJSON []byte
		if err := rows.Scan(&defJSON); err != nil {
			return nil, fmt.Errorf("failed to scan definition: %w", err)
		}

		var def domain.WorkflowDefinition
		if err := json.Unmarshal(defJSON, &def); err != nil {
			return nil, fmt.Errorf("failed to unmarshal definition: %w", err)
		}

		definitions = append(definitions, &def)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating definitions: %w", err)
	}

	return definitions, nil
}

// Delete removes a workflow definition by ID and version
func (r *DefinitionRepo) Delete(ctx context.Context, id string, version int) error {
	query := `DELETE FROM workflow_definitions WHERE id = $1 AND version = $2`

	result, err := r.pool.Exec(ctx, query, id, version)
	if err != nil {
		return fmt.Errorf("failed to delete definition: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("definition not found: %s v%d", id, version)
	}

	return nil
}
