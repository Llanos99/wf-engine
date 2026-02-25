package service

import (
	"context"
	"errors"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store"
)

// WorkflowService handles workflow definition operations
type WorkflowService struct {
	store store.Store
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(store store.Store) *WorkflowService {
	return &WorkflowService{store: store}
}

// CreateDefinition creates a new workflow definition
func (s *WorkflowService) CreateDefinition(ctx context.Context, def *domain.WorkflowDefinition) error {
	// Set default version
	if def.Version == 0 {
		def.Version = 1
	}

	// Validate
	if err := def.Validate(); err != nil {
		return err
	}

	// Save
	if err := s.store.Definitions().Save(ctx, def); err != nil {
		return &domain.DefinitionError{
			ID:      def.ID,
			Version: def.Version,
			Err:     err,
			Message: "failed to save",
		}
	}

	return nil
}

// GetDefinition retrieves a specific version of a definition
func (s *WorkflowService) GetDefinition(ctx context.Context, id string, version int) (*domain.WorkflowDefinition, error) {
	def, err := s.store.Definitions().Get(ctx, id, version)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewDefinitionNotFound(id, version)
		}
		return nil, err
	}
	return def, nil
}

// GetLatestDefinition retrieves the latest version of a definition
func (s *WorkflowService) GetLatestDefinition(ctx context.Context, id string) (*domain.WorkflowDefinition, error) {
	def, err := s.store.Definitions().GetLatest(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewDefinitionNotFound(id, 0)
		}
		return nil, err
	}
	return def, nil
}

// ListDefinitions retrieves all definitions
func (s *WorkflowService) ListDefinitions(ctx context.Context) ([]*domain.WorkflowDefinition, error) {
	return s.store.Definitions().List(ctx)
}

// DeleteDefinition deletes a specific version of a definition
func (s *WorkflowService) DeleteDefinition(ctx context.Context, id string, version int) error {
	// Check if any instances are using this definition
	instances, err := s.store.Instances().List(ctx, store.InstanceFilter{
		DefinitionID: id,
		Limit:        1,
	})
	if err != nil {
		return err
	}

	// Only warn - don't block deletion
	if len(instances) > 0 {
		// Could emit an event here
	}

	if err := s.store.Definitions().Delete(ctx, id, version); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.NewDefinitionNotFound(id, version)
		}
		return err
	}

	return nil
}
