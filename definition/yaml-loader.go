package definition

import (
	"os"

	"github.com/Llanos99/wf-engine/internal/domain"
	"gopkg.in/yaml.v3"
)

type YAMLLoader struct{}

func (l *YAMLLoader) Load(path string) (wf *domain.WorkflowDefinition, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var generatedWF domain.WorkflowDefinition
	if err := yaml.Unmarshal(data, &generatedWF); err != nil {
		return nil, err
	}
	return &generatedWF, nil
}
