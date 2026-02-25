package domain

import "time"

// ApprovalStatus represents the status of an approval request
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
	ApprovalExpired  ApprovalStatus = "expired"
)

// ApprovalRequest represents a pending approval in a workflow
type ApprovalRequest struct {
	ID          string         `json:"id"`
	InstanceID  string         `json:"instance_id"`
	NodeID      string         `json:"node_id"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Approvers   []string       `json:"approvers,omitempty"` // Allowed approvers
	Status      ApprovalStatus `json:"status"`
	DecidedBy   string         `json:"decided_by,omitempty"`
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	DecidedAt   *time.Time     `json:"decided_at,omitempty"`
}

// NewApprovalRequest creates a new pending approval request
func NewApprovalRequest(id, instanceID, nodeID, title string) *ApprovalRequest {
	return &ApprovalRequest{
		ID:         id,
		InstanceID: instanceID,
		NodeID:     nodeID,
		Title:      title,
		Status:     ApprovalPending,
		CreatedAt:  time.Now(),
	}
}

// Approve marks the request as approved
func (a *ApprovalRequest) Approve(decidedBy, notes string) {
	now := time.Now()
	a.Status = ApprovalApproved
	a.DecidedBy = decidedBy
	a.Notes = notes
	a.DecidedAt = &now
}

// Reject marks the request as rejected
func (a *ApprovalRequest) Reject(decidedBy, notes string) {
	now := time.Now()
	a.Status = ApprovalRejected
	a.DecidedBy = decidedBy
	a.Notes = notes
	a.DecidedAt = &now
}

// IsExpired checks if the approval has expired
func (a *ApprovalRequest) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsPending checks if the approval is still pending
func (a *ApprovalRequest) IsPending() bool {
	return a.Status == ApprovalPending && !a.IsExpired()
}
