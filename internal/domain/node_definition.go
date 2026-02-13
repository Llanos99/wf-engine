package domain

import "encoding/json"

// NodeType representa los tipos de nodos disponibles en el workflow
type NodeType string

const (
	NodeTypeStart       NodeType = "start"
	NodeTypeEnd         NodeType = "end"
	NodeTypeRunScript   NodeType = "run_script"
	NodeTypeSetValues   NodeType = "set_values"
	NodeTypeConditional NodeType = "conditional"
	NodeTypeWait        NodeType = "wait"
	NodeTypeLog         NodeType = "log"
)

// NodeDefinition es la definición de un nodo en el workflow
type NodeDefinition struct {
	ID   string          `yaml:"id" json:"id"`
	Type NodeType        `yaml:"type" json:"type"`
	Spec json.RawMessage `yaml:"spec" json:"spec"`
}
