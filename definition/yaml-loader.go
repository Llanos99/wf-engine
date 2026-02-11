package definition

import (
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
	var wf domain.WorkflowDefinition
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}
	return &wf, nil
}
