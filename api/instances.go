package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/store"
	"github.com/google/uuid"
)

// CreateInstanceRequest is the request body for creating an instance
type CreateInstanceRequest struct {
	DefinitionID string           `json:"definition_id"`
	Version      int              `json:"version,omitempty"` // Optional, defaults to latest
	Variables    domain.Variables `json:"variables,omitempty"`
}

// InstanceResponse is the response for an instance
type InstanceResponse struct {
	ID           string               `json:"id"`
	DefinitionID string               `json:"definition_id"`
	Version      int                  `json:"version"`
	Status       domain.WorkflowStatus `json:"status"`
	CurrentNode  string               `json:"current_node"`
	Variables    domain.Variables     `json:"variables"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	FinishedAt   *time.Time           `json:"finished_at,omitempty"`
}

// InstanceListResponse is a summary response for listing instances
type InstanceListResponse struct {
	ID           string               `json:"id"`
	DefinitionID string               `json:"definition_id"`
	Version      int                  `json:"version"`
	Status       domain.WorkflowStatus `json:"status"`
	CurrentNode  string               `json:"current_node"`
	CreatedAt    time.Time            `json:"created_at"`
}

// PendingApprovalInfo is included when instance is waiting for approval
type PendingApprovalInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// InstanceDetailResponse includes pending approval info
type InstanceDetailResponse struct {
	InstanceResponse
	PendingApproval *PendingApprovalInfo `json:"pending_approval,omitempty"`
	Logs            []domain.LogEntry    `json:"logs,omitempty"`
}

// instanceToResponse converts an instance to a response
func instanceToResponse(inst *domain.WorkflowInstance) InstanceResponse {
	return InstanceResponse{
		ID:           inst.ID,
		DefinitionID: inst.DefinitionID,
		Version:      inst.Version,
		Status:       inst.Status,
		CurrentNode:  inst.CurrentNode,
		Variables:    inst.Variables,
		CreatedAt:    inst.CreatedAt,
		UpdatedAt:    inst.UpdatedAt,
		FinishedAt:   inst.FinishedAt,
	}
}

// listInstances handles GET /api/v1/instances
func (s *Server) listInstances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := store.InstanceFilter{
		DefinitionID: r.URL.Query().Get("definition_id"),
		Status:       domain.WorkflowStatus(r.URL.Query().Get("status")),
		Limit:        100,
	}

	instances, err := s.store.Instances().List(ctx, filter)
	if err != nil {
		InternalError(w, err.Error())
		return
	}

	response := make([]InstanceListResponse, len(instances))
	for i, inst := range instances {
		response[i] = InstanceListResponse{
			ID:           inst.ID,
			DefinitionID: inst.DefinitionID,
			Version:      inst.Version,
			Status:       inst.Status,
			CurrentNode:  inst.CurrentNode,
			CreatedAt:    inst.CreatedAt,
		}
	}

	JSON(w, http.StatusOK, response)
}

// createInstance handles POST /api/v1/instances
func (s *Server) createInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if req.DefinitionID == "" {
		BadRequest(w, "definition_id is required")
		return
	}

	// Get the definition
	var def *domain.WorkflowDefinition
	var err error
	if req.Version > 0 {
		def, err = s.store.Definitions().Get(ctx, req.DefinitionID, req.Version)
	} else {
		def, err = s.store.Definitions().GetLatest(ctx, req.DefinitionID)
	}
	if err != nil {
		NotFound(w, "definition not found: "+req.DefinitionID)
		return
	}

	// Create instance
	instanceID := uuid.New().String()
	instance := domain.NewWorkflowInstance(instanceID, def)

	// Override variables if provided
	for k, v := range req.Variables {
		instance.Variables[k] = v
	}

	// Save instance
	if err := s.store.Instances().Create(ctx, instance); err != nil {
		InternalError(w, err.Error())
		return
	}

	// Execute workflow
	output, err := s.executor.Execute(def, instance)
	if err != nil {
		// Update instance with failed status
		s.store.Instances().Update(ctx, instance)
		InternalError(w, err.Error())
		return
	}

	// Handle pending approval
	if output != nil && output.PendingApproval != nil {
		if err := s.store.Approvals().Create(ctx, output.PendingApproval); err != nil {
			InternalError(w, "failed to create approval: "+err.Error())
			return
		}
	}

	// Update instance state
	if err := s.store.Instances().Update(ctx, instance); err != nil {
		InternalError(w, err.Error())
		return
	}

	// Build response
	response := InstanceDetailResponse{
		InstanceResponse: instanceToResponse(instance),
		Logs:             instance.Logs,
	}

	if output != nil && output.PendingApproval != nil {
		response.PendingApproval = &PendingApprovalInfo{
			ID:          output.PendingApproval.ID,
			Title:       output.PendingApproval.Title,
			Description: output.PendingApproval.Description,
		}
	}

	JSON(w, http.StatusCreated, response)
}

// getInstance handles GET /api/v1/instances/{id}
func (s *Server) getInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	instance, err := s.store.Instances().Get(ctx, id)
	if err != nil {
		NotFound(w, "instance not found: "+id)
		return
	}

	response := InstanceDetailResponse{
		InstanceResponse: instanceToResponse(instance),
		Logs:             instance.Logs,
	}

	// Check for pending approval
	if instance.Status == domain.StatusWaiting {
		approval, err := s.store.Approvals().GetPendingByInstance(ctx, id)
		if err == nil && approval != nil {
			response.PendingApproval = &PendingApprovalInfo{
				ID:          approval.ID,
				Title:       approval.Title,
				Description: approval.Description,
			}
		}
	}

	JSON(w, http.StatusOK, response)
}

// deleteInstance handles DELETE /api/v1/instances/{id}
func (s *Server) deleteInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	if err := s.store.Instances().Delete(ctx, id); err != nil {
		NotFound(w, "instance not found: "+id)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// resumeInstance handles POST /api/v1/instances/{id}/resume
func (s *Server) resumeInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	instance, err := s.store.Instances().Get(ctx, id)
	if err != nil {
		NotFound(w, "instance not found: "+id)
		return
	}

	if instance.Status != domain.StatusWaiting {
		Conflict(w, "instance is not in waiting status")
		return
	}

	// Get the definition
	def, err := s.store.Definitions().Get(ctx, instance.DefinitionID, instance.Version)
	if err != nil {
		InternalError(w, "definition not found")
		return
	}

	// Check if waiting for approval
	approval, _ := s.store.Approvals().GetPendingByInstance(ctx, id)
	if approval != nil && approval.IsPending() {
		Conflict(w, "instance is waiting for approval: "+approval.ID)
		return
	}

	// Resume execution
	output, err := s.executor.Resume(def, instance)
	if err != nil {
		s.store.Instances().Update(ctx, instance)
		InternalError(w, err.Error())
		return
	}

	// Handle new pending approval
	if output != nil && output.PendingApproval != nil {
		if err := s.store.Approvals().Create(ctx, output.PendingApproval); err != nil {
			InternalError(w, "failed to create approval: "+err.Error())
			return
		}
	}

	// Update instance
	if err := s.store.Instances().Update(ctx, instance); err != nil {
		InternalError(w, err.Error())
		return
	}

	response := InstanceDetailResponse{
		InstanceResponse: instanceToResponse(instance),
		Logs:             instance.Logs,
	}

	if output != nil && output.PendingApproval != nil {
		response.PendingApproval = &PendingApprovalInfo{
			ID:          output.PendingApproval.ID,
			Title:       output.PendingApproval.Title,
			Description: output.PendingApproval.Description,
		}
	}

	JSON(w, http.StatusOK, response)
}

// getInstanceLogs handles GET /api/v1/instances/{id}/logs
func (s *Server) getInstanceLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	logs, err := s.store.Instances().GetLogs(ctx, id)
	if err != nil {
		NotFound(w, "instance not found: "+id)
		return
	}

	JSON(w, http.StatusOK, logs)
}
