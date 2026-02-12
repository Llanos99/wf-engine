package runtime

import "time"

type ExecutionResult struct {
	NextStepID string
	WaitUntil  *time.Time
	Error      error
	Finished   bool
}
