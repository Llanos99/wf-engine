package definition

import "github.com/Llanos99/wf-engine/internal/domain"

type Loader interface {
	Load(path string) (*domain.WorkflowDefinition, error)
}

var LoaderTypes = map[string]Loader{
	"yaml": &YAMLLoader{},
}
