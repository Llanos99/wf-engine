package definition

import (
	"errors"

	"github.com/Llanos99/wf-engine/internal/domain"
)

type Loader interface {
	Load(path string) (wf *domain.WorkflowDefinition, err error)
}

var LoaderTypes = map[string]Loader{
	"yaml": &YAMLLoader{},
}

func GetLoader(loaderType string) (Loader, error) {
	loader, ok := LoaderTypes[loaderType]
	if !ok {
		return nil, errors.New("Unknown loader type")
	}
	return loader, nil
}
