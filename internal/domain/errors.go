package domain

import (
	"errors"
	"fmt"
)

// Base errors
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidState  = errors.New("invalid state")
	ErrValidation    = errors.New("validation error")
)

// DefinitionError represents an error related to workflow definitions
type DefinitionError struct {
	ID      string
	Version int
	Err     error
	Message string
}

func (e *DefinitionError) Error() string {
	if e.Version > 0 {
		return fmt.Sprintf("definition %s v%d: %s", e.ID, e.Version, e.Message)
	}
	return fmt.Sprintf("definition %s: %s", e.ID, e.Message)
}

func (e *DefinitionError) Unwrap() error {
	return e.Err
}

// InstanceError represents an error related to workflow instances
type InstanceError struct {
	ID      string
	Err     error
	Message string
}

func (e *InstanceError) Error() string {
	return fmt.Sprintf("instance %s: %s", e.ID, e.Message)
}

func (e *InstanceError) Unwrap() error {
	return e.Err
}

// ApprovalError represents an error related to approvals
type ApprovalError struct {
	ID      string
	Err     error
	Message string
}

func (e *ApprovalError) Error() string {
	return fmt.Sprintf("approval %s: %s", e.ID, e.Message)
}

func (e *ApprovalError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}
	return fmt.Sprintf("%d validation errors", len(e))
}

func (e ValidationErrors) Unwrap() error {
	return ErrValidation
}

// Helper constructors
func NewDefinitionNotFound(id string, version int) error {
	return &DefinitionError{ID: id, Version: version, Err: ErrNotFound, Message: "not found"}
}

func NewInstanceNotFound(id string) error {
	return &InstanceError{ID: id, Err: ErrNotFound, Message: "not found"}
}

func NewApprovalNotFound(id string) error {
	return &ApprovalError{ID: id, Err: ErrNotFound, Message: "not found"}
}

func NewInstanceInvalidState(id string, current, expected WorkflowStatus) error {
	return &InstanceError{
		ID:      id,
		Err:     ErrInvalidState,
		Message: fmt.Sprintf("status is '%s', expected '%s'", current, expected),
	}
}

func NewApprovalAlreadyDecided(id string, status ApprovalStatus) error {
	return &ApprovalError{
		ID:      id,
		Err:     ErrInvalidState,
		Message: fmt.Sprintf("already decided: %s", status),
	}
}

func NewApprovalExpired(id string) error {
	return &ApprovalError{
		ID:      id,
		Err:     ErrInvalidState,
		Message: "approval has expired",
	}
}
