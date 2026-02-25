package nodes

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// ApprovalHandler maneja nodos de tipo "approval"
type ApprovalHandler struct{}

func NewApprovalHandler() *ApprovalHandler {
	return &ApprovalHandler{}
}

func (h *ApprovalHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.ApprovalSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid approval spec: %w", err)
	}

	if spec.Title == "" {
		return nil, fmt.Errorf("approval node '%s' has no title", n.ID)
	}

	if spec.ApprovedNext == "" {
		return nil, fmt.Errorf("approval node '%s' has no approved_next", n.ID)
	}

	// Interpolate variables in title and description
	title := h.interpolate(spec.Title, ctx.Instance.Variables)
	description := h.interpolate(spec.Description, ctx.Instance.Variables)

	// Generate approval request ID
	approvalID := fmt.Sprintf("%s-%s-%d", ctx.Instance.ID, n.ID, time.Now().UnixNano())

	// Create approval request
	approval := domain.NewApprovalRequest(approvalID, ctx.Instance.ID, n.ID, title)
	approval.Description = description
	approval.Approvers = spec.Approvers

	// Parse timeout if specified
	if spec.Timeout != "" {
		duration, err := parseDuration(spec.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout '%s': %w", spec.Timeout, err)
		}
		expiresAt := time.Now().Add(duration)
		approval.ExpiresAt = &expiresAt
	}

	ctx.Instance.AddLog("info", n.ID, fmt.Sprintf("Approval requested: %s", title))

	// Return with approval request - executor will persist it
	return &node.ExecutionResult{
		WaitForApproval: approval,
	}, nil
}

// interpolate replaces ${var} with values from variables
func (h *ApprovalHandler) interpolate(s string, vars domain.Variables) string {
	result := s
	for key, value := range vars {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// parseDuration parses durations like "24h", "7d", "30m"
func parseDuration(s string) (time.Duration, error) {
	// Handle days specially since Go doesn't support "d"
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		var d int
		if _, err := fmt.Sscanf(days, "%d", &d); err != nil {
			return 0, err
		}
		return time.Duration(d) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
