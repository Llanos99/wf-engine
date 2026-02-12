package definition

import (
	"encoding/json"
	"os"

	"github.com/Llanos99/wf-engine/internal/domain"
	"gopkg.in/yaml.v3"
)

type YAMLLoader struct{}

func (l *YAMLLoader) Load(path string) (*domain.WorkflowDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var def domain.WorkflowDefinition
	if err := json.Unmarshal(jsonData, &def); err != nil {
		return nil, err
	}
	return &def, nil
}
