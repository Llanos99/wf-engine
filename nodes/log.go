package nodes

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/node"
)

// LogHandler maneja nodos de tipo "log"
type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) Execute(ctx *node.ExecutionContext, n domain.NodeDefinition) (*node.ExecutionResult, error) {
	var spec domain.LogSpec
	if err := json.Unmarshal(n.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid log spec: %w", err)
	}

	// Interpolar variables en el mensaje
	message := h.interpolate(spec.Message, ctx.Instance.Variables)

	// Determinar nivel de log
	level := spec.Level
	if level == "" {
		level = "info"
	}

	// Registrar el log
	ctx.Instance.AddLog(level, n.ID, message)
	fmt.Printf("[%s] %s\n", level, message)

	return &node.ExecutionResult{
		NextNodeID: spec.Next,
	}, nil
}

// interpolate reemplaza ${variable} con los valores de las variables
func (h *LogHandler) interpolate(message string, vars domain.Variables) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(message, func(match string) string {
		// Extraer el nombre de la variable (sin ${ y })
		varName := match[2 : len(match)-1]
		if val, ok := vars.Get(varName); ok {
			return fmt.Sprintf("%v", val)
		}
		return match // Mantener el original si no existe
	})
}
