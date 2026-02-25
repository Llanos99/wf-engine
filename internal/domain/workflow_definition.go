package domain

import "encoding/json"

// WorkflowDefinition es el blueprint inmutable de un workflow
type WorkflowDefinition struct {
	ID               string                    `yaml:"id" json:"id"`
	Version          int                       `yaml:"version" json:"version"`
	Name             string                    `yaml:"name" json:"name"`
	Description      string                    `yaml:"description" json:"description"`
	Nodes            map[string]NodeDefinition `yaml:"nodes" json:"nodes"`
	InitialVariables Variables                 `yaml:"initial_variables" json:"initial_variables"`
}

// GetStartNode encuentra el nodo de tipo "start"
func (w *WorkflowDefinition) GetStartNode() (*NodeDefinition, bool) {
	for _, node := range w.Nodes {
		if node.Type == NodeTypeStart {
			return &node, true
		}
	}
	return nil, false
}

// GetNode obtiene un nodo por su ID
func (w *WorkflowDefinition) GetNode(id string) (*NodeDefinition, bool) {
	node, ok := w.Nodes[id]
	if !ok {
		return nil, false
	}
	return &node, true
}

// Validate verifica que el workflow sea válido
func (w *WorkflowDefinition) Validate() error {
	var errs ValidationErrors

	// Basic field validation
	if w.ID == "" {
		errs = append(errs, ValidationError{Field: "id", Message: "is required"})
	}
	if w.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}
	if len(w.Nodes) == 0 {
		errs = append(errs, ValidationError{Field: "nodes", Message: "at least one node is required"})
		return errs
	}

	// Count start and end nodes
	var startNodes, endNodes []string
	for id, node := range w.Nodes {
		if node.Type == NodeTypeStart {
			startNodes = append(startNodes, id)
		}
		if node.Type == NodeTypeEnd {
			endNodes = append(endNodes, id)
		}
	}

	// Must have exactly one start node
	if len(startNodes) == 0 {
		errs = append(errs, ValidationError{Field: "nodes", Message: "must have exactly one start node"})
	} else if len(startNodes) > 1 {
		errs = append(errs, ValidationError{Field: "nodes", Message: "must have exactly one start node, found multiple"})
	}

	// Must have at least one end node
	if len(endNodes) == 0 {
		errs = append(errs, ValidationError{Field: "nodes", Message: "must have at least one end node"})
	}

	// Validate node references
	reachable := make(map[string]bool)
	if len(startNodes) == 1 {
		w.markReachable(startNodes[0], reachable)
	}

	for id, node := range w.Nodes {
		// Check that node ID matches map key
		if node.ID != id {
			errs = append(errs, ValidationError{
				Field:   "nodes." + id,
				Message: "node ID does not match key",
			})
		}

		// Validate next node references
		nextNodes := w.getNextNodes(node)
		for _, next := range nextNodes {
			if next != "" {
				if _, exists := w.Nodes[next]; !exists {
					errs = append(errs, ValidationError{
						Field:   "nodes." + id,
						Message: "references non-existent node: " + next,
					})
				}
			}
		}
	}

	// Check for orphan nodes (unreachable from start)
	for id := range w.Nodes {
		if !reachable[id] && len(startNodes) == 1 {
			errs = append(errs, ValidationError{
				Field:   "nodes." + id,
				Message: "node is unreachable from start",
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// markReachable marks all nodes reachable from the given node
func (w *WorkflowDefinition) markReachable(nodeID string, reachable map[string]bool) {
	if reachable[nodeID] {
		return // Already visited
	}
	reachable[nodeID] = true

	node, ok := w.Nodes[nodeID]
	if !ok {
		return
	}

	for _, next := range w.getNextNodes(node) {
		if next != "" {
			w.markReachable(next, reachable)
		}
	}
}

// getNextNodes extracts all possible next node IDs from a node's spec
func (w *WorkflowDefinition) getNextNodes(node NodeDefinition) []string {
	var result []string

	// Parse spec based on node type
	switch node.Type {
	case NodeTypeStart, NodeTypeRunScript, NodeTypeSetValues, NodeTypeWait, NodeTypeLog:
		var spec struct {
			Next string `json:"next"`
		}
		if json.Unmarshal(node.Spec, &spec) == nil {
			result = append(result, spec.Next)
		}
	case NodeTypeConditional:
		var spec struct {
			TrueNext  string `json:"true_next"`
			FalseNext string `json:"false_next"`
		}
		if json.Unmarshal(node.Spec, &spec) == nil {
			result = append(result, spec.TrueNext, spec.FalseNext)
		}
	case NodeTypeApproval:
		var spec struct {
			ApprovedNext string `json:"approved_next"`
			RejectedNext string `json:"rejected_next"`
		}
		if json.Unmarshal(node.Spec, &spec) == nil {
			result = append(result, spec.ApprovedNext, spec.RejectedNext)
		}
	case NodeTypeEnd:
		// End nodes have no next
	}

	return result
}
