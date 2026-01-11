package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Llanos99/wf-engine/engine"
	"github.com/Llanos99/wf-engine/models"
)

func main() {
	ctx := &models.Context{
		Variables:      &models.Variables{Data: make(map[string]models.Variable)},
		StepResults:    make(map[string]interface{}),
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		ExecutionCount: make(map[string]int),
	}

	// Populate variables
	ctx.Variables.SetInt("variable_1", 1)

	step1 := models.Step{
		ID:   "first_step",
		Name: "First step",
		Type: "if-else",
		Config: map[string]interface{}{
			"condition": func(ctx *models.Context) bool {
				ctx.Logger.Println("Loading first step")
				v, err := ctx.Variables.GetInt("variable_1")
				if err != nil {
					fmt.Printf("variable %s corrupted", "variable_1")
				}
				return v > 0
			},
			"true_next":  "second_step",
			"false_next": "third_step",
		},
		NextID: "second_step",
	}

	step2 := models.Step{
		ID:   "second_step",
		Name: "Second step",
		Type: "if-else",
		Config: map[string]interface{}{
			"condition": func(ctx *models.Context) bool {
				ctx.Logger.Println("Loading second step")
				v, err := ctx.Variables.GetInt("variable_1")
				if err != nil {
					fmt.Printf("variable %s corrupted", "variable_1")
				}
				return v > 0
			},
			"true_next":  "first_step",
			"false_next": "third_step",
		},
		NextID: "third_step",
	}

	step3 := models.Step{
		ID:     "third_step",
		Name:   "Third step",
		Type:   "action",
		Config: map[string]interface{}{},
		NextID: "",
	}

	wf := &models.Workflow{
		StartAt: "first_step",
		Steps: map[string]*models.Step{
			"first_step":  &step1,
			"second_step": &step2,
			"third_step":  &step3,
		},
	}

	executor := &engine.Executor{}

	err := executor.Run(wf, ctx)
	if err != nil {
		fmt.Println("Workflow error: ", err)
		return
	}
	fmt.Println("Workflow completed!")

	// What's next?
}
