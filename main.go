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
		Data:           make(map[string]interface{}),
		StepResults:    make(map[string]interface{}),
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		ExecutionCount: make(map[string]int),
	}

	// Populate variables
	ctx.Data["variable_1"] = map[string]interface{}{
		"val": 0,
	}

	step1 := models.Step{
		ID:   "first_step",
		Name: "First step",
		Type: "if-else",
		Config: map[string]interface{}{
			"condition": func(ctx *models.Context) bool {
				ctx.Logger.Println("Loading first step")
				v := ctx.Data["variable_1"].(map[string]interface{})["val"].(int)
				return v > 0
			},
			"true_next":  "second_step",
			"false_next": "third_step",
		},
		NextID: "second_step",
	}

	step2 := models.Step{
		ID:     "second_step",
		Name:   "Second step",
		Type:   "action",
		Config: map[string]interface{}{},
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
