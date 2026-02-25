package domain

// StartSpec define el nodo de inicio del workflow
type StartSpec struct {
	Next string `yaml:"next" json:"next"`
}

// EndSpec define el nodo de fin del workflow
type EndSpec struct {
	// El nodo End no necesita configuración adicional
}

// RunScriptSpec define un nodo que ejecuta JavaScript
type RunScriptSpec struct {
	Script string `yaml:"script" json:"script"` // Código JavaScript a ejecutar
	Next   string `yaml:"next" json:"next"`
}

// SetValuesSpec define un nodo que asigna valores a variables
type SetValuesSpec struct {
	Assignments map[string]interface{} `yaml:"assignments" json:"assignments"`
	Next        string                 `yaml:"next" json:"next"`
}

// ConditionalSpec define un nodo de branching condicional
type ConditionalSpec struct {
	// Puede usar Expression o Script para evaluar la condición
	Expression *Expression `yaml:"expression,omitempty" json:"expression,omitempty"`
	Script     string      `yaml:"script,omitempty" json:"script,omitempty"` // JS que retorna boolean
	TrueNext   string      `yaml:"true_next" json:"true_next"`
	FalseNext  string      `yaml:"false_next" json:"false_next"`
}

// WaitSpec define un nodo que pausa la ejecución
type WaitSpec struct {
	Duration string `yaml:"duration" json:"duration"` // Go duration: "10s", "5m", "1h"
	Next     string `yaml:"next" json:"next"`
}

// LogSpec define un nodo que registra información
type LogSpec struct {
	Level   string `yaml:"level" json:"level"`     // debug, info, warn, error
	Message string `yaml:"message" json:"message"` // Puede contener ${variable} para interpolación
	Next    string `yaml:"next" json:"next"`
}

// ApprovalSpec define un nodo que requiere aprobación humana
type ApprovalSpec struct {
	Title        string   `yaml:"title" json:"title"`                                   // Título de la aprobación
	Description  string   `yaml:"description,omitempty" json:"description,omitempty"`   // Descripción detallada
	Approvers    []string `yaml:"approvers,omitempty" json:"approvers,omitempty"`       // Lista de aprobadores permitidos
	Timeout      string   `yaml:"timeout,omitempty" json:"timeout,omitempty"`           // Tiempo límite: "24h", "7d"
	ApprovedNext string   `yaml:"approved_next" json:"approved_next"`                   // Siguiente nodo si aprobado
	RejectedNext string   `yaml:"rejected_next,omitempty" json:"rejected_next,omitempty"` // Siguiente nodo si rechazado (opcional)
}
