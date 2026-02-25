package store

import (
	"context"

	"github.com/Llanos99/wf-engine/internal/domain"
)

// DefinitionRepository maneja la persistencia de WorkflowDefinitions
type DefinitionRepository interface {
	// Save guarda una definición de workflow
	Save(ctx context.Context, def *domain.WorkflowDefinition) error

	// Get obtiene una definición por ID y versión
	Get(ctx context.Context, id string, version int) (*domain.WorkflowDefinition, error)

	// GetLatest obtiene la última versión de una definición
	GetLatest(ctx context.Context, id string) (*domain.WorkflowDefinition, error)

	// List lista todas las definiciones (última versión de cada una)
	List(ctx context.Context) ([]*domain.WorkflowDefinition, error)

	// Delete elimina una definición
	Delete(ctx context.Context, id string, version int) error
}

// InstanceRepository maneja la persistencia de WorkflowInstances
type InstanceRepository interface {
	// Create crea una nueva instancia
	Create(ctx context.Context, instance *domain.WorkflowInstance) error

	// Update actualiza una instancia existente
	Update(ctx context.Context, instance *domain.WorkflowInstance) error

	// Get obtiene una instancia por ID
	Get(ctx context.Context, id string) (*domain.WorkflowInstance, error)

	// List lista instancias con filtros opcionales
	List(ctx context.Context, filter InstanceFilter) ([]*domain.WorkflowInstance, error)

	// Delete elimina una instancia
	Delete(ctx context.Context, id string) error

	// SaveLog guarda una entrada de log
	SaveLog(ctx context.Context, instanceID string, log domain.LogEntry) error

	// GetLogs obtiene los logs de una instancia
	GetLogs(ctx context.Context, instanceID string) ([]domain.LogEntry, error)
}

// InstanceFilter define filtros para listar instancias
type InstanceFilter struct {
	DefinitionID string
	Status       domain.WorkflowStatus
	Limit        int
	Offset       int
}

// ApprovalRepository maneja la persistencia de ApprovalRequests
type ApprovalRepository interface {
	// Create crea una nueva solicitud de aprobación
	Create(ctx context.Context, approval *domain.ApprovalRequest) error

	// Get obtiene una solicitud por ID
	Get(ctx context.Context, id string) (*domain.ApprovalRequest, error)

	// GetByInstance obtiene todas las solicitudes de una instancia
	GetByInstance(ctx context.Context, instanceID string) ([]*domain.ApprovalRequest, error)

	// GetPending obtiene la solicitud pendiente de una instancia (si existe)
	GetPendingByInstance(ctx context.Context, instanceID string) (*domain.ApprovalRequest, error)

	// Update actualiza una solicitud (para aprobar/rechazar)
	Update(ctx context.Context, approval *domain.ApprovalRequest) error

	// ListPending lista todas las solicitudes pendientes
	ListPending(ctx context.Context) ([]*domain.ApprovalRequest, error)

	// Delete elimina una solicitud
	Delete(ctx context.Context, id string) error
}

// Store agrupa todos los repositorios
type Store interface {
	Definitions() DefinitionRepository
	Instances() InstanceRepository
	Approvals() ApprovalRepository
	Close() error
}
