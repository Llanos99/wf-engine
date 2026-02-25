package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Llanos99/wf-engine/definition"
	"github.com/Llanos99/wf-engine/engine"
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/nodes"
	"github.com/Llanos99/wf-engine/script"
)

func main() {
	// Determinar qué workflow cargar
	workflowPath := "testdata/example_workflow.yaml"
	if len(os.Args) > 1 {
		workflowPath = os.Args[1]
	}

	fmt.Println("=== WF-Engine ===")
	fmt.Printf("Loading workflow: %s\n\n", workflowPath)

	// Cargar definición del workflow
	loader, err := definition.GetLoader("yaml")
	if err != nil {
		panic(err)
	}

	def, err := loader.Load(workflowPath)
	if err != nil {
		panic(fmt.Errorf("failed to load workflow: %w", err))
	}

	fmt.Printf("Workflow: %s (v%d)\n", def.Name, def.Version)
	fmt.Printf("Description: %s\n", def.Description)
	fmt.Printf("Nodes: %d\n\n", len(def.Nodes))

	// Crear instancia del workflow
	instance := domain.NewWorkflowInstance("instance-001", def)

	// Crear registry de nodos con handlers builtin
	nodeRegistry := nodes.NewRegistryWithBuiltin()

	// Crear script runner
	scriptRunner := script.NewAdapter()

	// Crear executor
	executor := engine.NewExecutor(nodeRegistry, scriptRunner)

	fmt.Println("--- Execution Start ---")

	// Ejecutar
	output, err := executor.Execute(def, instance)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Execution End ---")

	// Check if waiting for approval
	if output != nil && output.PendingApproval != nil {
		fmt.Printf("\nPending Approval: %s\n", output.PendingApproval.Title)
		fmt.Printf("  ID: %s\n", output.PendingApproval.ID)
		if output.PendingApproval.Description != "" {
			fmt.Printf("  Description: %s\n", output.PendingApproval.Description)
		}
	}

	// Mostrar resultado
	fmt.Printf("\nStatus: %s\n", instance.Status)
	fmt.Println("\nFinal Variables:")
	varsJSON, _ := json.MarshalIndent(instance.Variables, "", "  ")
	fmt.Println(string(varsJSON))

	fmt.Println("\nExecution Logs:")
	for _, log := range instance.Logs {
		fmt.Printf("  [%s] %s: %s\n", log.Level, log.NodeID, log.Message)
	}
}
