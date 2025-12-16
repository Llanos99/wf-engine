package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Llanos99/wf-engine/models"
)

func main() {
	ctx := &models.Context{
		Data:        make(map[string]interface{}),
		StepResults: make(map[string]interface{}),
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
	}

	step1 := models.Step{
		ID:   "first_step",
		Name: "First step",
		Type: "if",
		Execute: func(ctx *models.Context) error {
			ctx.Logger.Println("Loading first step")
			ctx.Data["variable_1"] = map[string]interface{}{
				"val": 1,
			}
			return nil
		},
		NextID: "second_step",
	}

	step2 := models.Step{
		ID:   "second_step",
		Name: "Second step",
		Type: "if",
		Execute: func(ctx *models.Context) error {
			ctx.Logger.Println("Loading second step")
			var_1 := ctx.Data["variable_1"].(map[string]interface{})
			val_1 := var_1["val"].(int)
			ctx.Logger.Println("Received data from first step: ", val_1)
			ctx.Data["variable_2"] = map[string]interface{}{
				"val": 2 * val_1,
			}
			ctx.StepResults["second_step"] = map[string]interface{}{
				"status": "done",
			}
			return nil
		},
		NextID: "third_step",
	}

	step3 := models.Step{
		ID:   "third_step",
		Name: "Third step",
		Type: "if",
		Execute: func(ctx *models.Context) error {
			ctx.Logger.Println("Loading third step")
			var_2 := ctx.Data["variable_2"].(map[string]interface{})
			val_2 := var_2["val"].(int)
			ctx.Logger.Println("Received data from first step: ", val_2)
			ctx.Data["variable_3"] = map[string]interface{}{
				"val": 2 * val_2,
			}
			ctx.StepResults["third_step"] = map[string]interface{}{
				"status": "done",
			}
			return nil
		},
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

	executor := &models.Executor{}

	err := executor.Run(wf, ctx)
	if err != nil {
		fmt.Println("Workflow error: ", err)
		return
	}
	fmt.Println("Workflow completed!")

	// What's next?
}
