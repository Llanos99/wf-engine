package strategy

import (
	"fmt"

	"github.com/Llanos99/wf-engine/models"
)

type ConditionalHandler struct{}

func (h *ConditionalHandler) Execute(ctx *models.Context, step *models.Step) (executionResult *models.ExecutionResult, err error) {
	ctx.Logger.Println("Executing IF step: ", step.ID)
	conditionalFunc, ok := step.Config["condition"].(func(*models.Context) bool)
	if !ok {
		return nil, fmt.Errorf("Conditional step %s has no condition", step.ID)
	}
	trueNext := step.Config["true_next"].(string)
	falseNext := step.Config["false_next"].(string)
	var result = conditionalFunc(ctx)
	// Write the StepResults
	ctx.StepResults[step.ID] = map[string]interface{}{
		"type":   "if",
		"result": result,
		"status": "done",
	}
	if result {
		ctx.Logger.Printf("Step %s condition TRUE then NEXT is %s", step.ID, trueNext)
		return &models.ExecutionResult{
			Status:   models.COMPLETED,
			NextStep: trueNext,
		}, nil
	}
	ctx.Logger.Printf("Step %s condition FALSE then NEXT is %s", step.ID, falseNext)
	return &models.ExecutionResult{
		Status:   models.COMPLETED,
		NextStep: falseNext,
	}, nil
}

func (h *ConditionalHandler) Validate(step *models.Step) error {
	if _, ok := step.Config["condition"].(func(*models.Context) bool); !ok {
		return fmt.Errorf("Step %s missing condition", step.ID)
	}
	if _, ok := step.Config["true_next"].(string); !ok {
		return fmt.Errorf("Step %s missing true_next", step.ID)
	}
	if _, ok := step.Config["false_next"].(string); !ok {
		return fmt.Errorf("Step %s missing false_next", step.ID)
	}
	return nil
}
