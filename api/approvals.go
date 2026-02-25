package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
)

// ApprovalResponse is the response for an approval
type ApprovalResponse struct {
	ID          string               `json:"id"`
	InstanceID  string               `json:"instance_id"`
	NodeID      string               `json:"node_id"`
	Title       string               `json:"title"`
	Description string               `json:"description,omitempty"`
	Approvers   []string             `json:"approvers,omitempty"`
	Status      domain.ApprovalStatus `json:"status"`
	DecidedBy   string               `json:"decided_by,omitempty"`
	Notes       string               `json:"notes,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	ExpiresAt   *time.Time           `json:"expires_at,omitempty"`
	DecidedAt   *time.Time           `json:"decided_at,omitempty"`
}

// ApprovalDecisionRequest is the request body for approving/rejecting
type ApprovalDecisionRequest struct {
	DecidedBy string `json:"decided_by"`
	Notes     string `json:"notes,omitempty"`
}

// ApprovalDecisionResponse includes the updated instance
type ApprovalDecisionResponse struct {
	Approval ApprovalResponse     `json:"approval"`
	Instance InstanceDetailResponse `json:"instance"`
}

// approvalToResponse converts an approval to a response
func approvalToResponse(a *domain.ApprovalRequest) ApprovalResponse {
	return ApprovalResponse{
		ID:          a.ID,
		InstanceID:  a.InstanceID,
		NodeID:      a.NodeID,
		Title:       a.Title,
		Description: a.Description,
		Approvers:   a.Approvers,
		Status:      a.Status,
		DecidedBy:   a.DecidedBy,
		Notes:       a.Notes,
		CreatedAt:   a.CreatedAt,
		ExpiresAt:   a.ExpiresAt,
		DecidedAt:   a.DecidedAt,
	}
}

// listPendingApprovals handles GET /api/v1/approvals
func (s *Server) listPendingApprovals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	approvals, err := s.store.Approvals().ListPending(ctx)
	if err != nil {
		InternalError(w, err.Error())
		return
	}

	response := make([]ApprovalResponse, len(approvals))
	for i, a := range approvals {
		response[i] = approvalToResponse(a)
	}

	JSON(w, http.StatusOK, response)
}

// getApproval handles GET /api/v1/approvals/{id}
func (s *Server) getApproval(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	approval, err := s.store.Approvals().Get(ctx, id)
	if err != nil {
		NotFound(w, "approval not found: "+id)
		return
	}

	JSON(w, http.StatusOK, approvalToResponse(approval))
}

// approveRequest handles POST /api/v1/approvals/{id}/approve
func (s *Server) approveRequest(w http.ResponseWriter, r *http.Request) {
	s.handleApprovalDecision(w, r, true)
}

// rejectRequest handles POST /api/v1/approvals/{id}/reject
func (s *Server) rejectRequest(w http.ResponseWriter, r *http.Request) {
	s.handleApprovalDecision(w, r, false)
}

// handleApprovalDecision processes approval or rejection
func (s *Server) handleApprovalDecision(w http.ResponseWriter, r *http.Request, approve bool) {
	ctx := r.Context()
	id := urlParam(r, "id")

	// Parse request
	var req ApprovalDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if req.DecidedBy == "" {
		BadRequest(w, "decided_by is required")
		return
	}

	// Get approval
	approval, err := s.store.Approvals().Get(ctx, id)
	if err != nil {
		NotFound(w, "approval not found: "+id)
		return
	}

	// Check if already decided
	if approval.Status != domain.ApprovalPending {
		Conflict(w, "approval already decided: "+string(approval.Status))
		return
	}

	// Check if expired
	if approval.IsExpired() {
		approval.Status = domain.ApprovalExpired
		s.store.Approvals().Update(ctx, approval)
		Conflict(w, "approval has expired")
		return
	}

	// Apply decision
	if approve {
		approval.Approve(req.DecidedBy, req.Notes)
	} else {
		approval.Reject(req.DecidedBy, req.Notes)
	}

	// Update approval
	if err := s.store.Approvals().Update(ctx, approval); err != nil {
		InternalError(w, err.Error())
		return
	}

	// Get instance
	instance, err := s.store.Instances().Get(ctx, approval.InstanceID)
	if err != nil {
		InternalError(w, "instance not found: "+approval.InstanceID)
		return
	}

	// Get definition
	def, err := s.store.Definitions().Get(ctx, instance.DefinitionID, instance.Version)
	if err != nil {
		InternalError(w, "definition not found")
		return
	}

	// Resume workflow after approval
	output, err := s.executor.ResumeAfterApproval(def, instance, approval)
	if err != nil {
		s.store.Instances().Update(ctx, instance)
		InternalError(w, err.Error())
		return
	}

	// Handle new pending approval (if workflow hit another approval node)
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

	// Build response
	instanceResponse := InstanceDetailResponse{
		InstanceResponse: instanceToResponse(instance),
		Logs:             instance.Logs,
	}

	if output != nil && output.PendingApproval != nil {
		instanceResponse.PendingApproval = &PendingApprovalInfo{
			ID:          output.PendingApproval.ID,
			Title:       output.PendingApproval.Title,
			Description: output.PendingApproval.Description,
		}
	}

	JSON(w, http.StatusOK, ApprovalDecisionResponse{
		Approval: approvalToResponse(approval),
		Instance: instanceResponse,
	})
}
