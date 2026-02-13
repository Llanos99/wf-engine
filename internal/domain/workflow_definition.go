package domain

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
	// TODO: Implementar validación
	// - Debe tener exactamente un nodo Start
	// - Debe tener al menos un nodo End
	// - Todos los Next deben apuntar a nodos existentes
	// - No debe haber nodos huérfanos
	return nil
}
