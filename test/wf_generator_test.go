package test

import (
	"encoding/json"
	"testing"

	"github.com/Llanos99/wf-engine/definition"
)

func TestWorkflowGenerator(t *testing.T) {
	path := "../testdata/simple.yaml"
	loaderType := "yaml"
	loader, err := definition.GetLoader(loaderType)
	if err != nil {
		t.Fatalf("Failed to get loader: %v", err)
	}
	wf, err := loader.Load(path)
	if err != nil {
		t.Fatalf("Failed to load workflow: %v", err)
	}
	pretty, err := json.MarshalIndent(wf, "", " ")
	if err != nil {
		t.Fatalf("Failed to marshal workflow: %v", err)
	}
	t.Logf("Workflow loaded successfully:\n%s", string(pretty))
}
