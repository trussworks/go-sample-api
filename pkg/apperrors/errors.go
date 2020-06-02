package apperrors

import (
	"fmt"

	"bin/bork/pkg/models"
)

// QueryOperation provides a set of operations that can fail
type QueryOperation string

const (
	// QueryPost is for failures when creating a resource
	QueryPost QueryOperation = "Create"
	// QuerySave is for failures when saving
	QuerySave QueryOperation = "Save"
	// QueryFetch is for failures when getting a resource
	QueryFetch QueryOperation = "Fetch"
)

// QueryError is a typed error for query issues
type QueryError struct {
	Err       error
	Resource  interface{}
	Operation QueryOperation
}

// Error provides the error as a string
func (e *QueryError) Error() string {
	return fmt.Sprintf(
		"Could not query model %T with operation %s, received error: %s",
		e.Resource,
		e.Operation,
		e.Err,
	)
}

// Unwrap provides the underlying error
func (e *QueryError) Unwrap() error {
	return e.Err
}

// UnauthorizedError is a typed error for when authorization fails
type UnauthorizedError struct {
	User      models.User
	Operation QueryOperation
	Resource  interface{}
	Err       error
}

// Error provides the error as a string
func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf(
		"User: %s is unauthorized for operation: %s on resource: %T with error: %s",
		e.User.ID,
		e.Operation,
		e.Resource,
		e.Err,
	)
}

// Unwrap provides the underlying error
func (e *UnauthorizedError) Unwrap() error {
	return e.Err
}

// ContextResource is an enum for context resources
type ContextResource string

const (
	// ContextResourceUser is a constant representing User on context
	ContextResourceUser ContextResource = "User"
	// ContextResourceLogger is a constant representing Logger on context
	ContextResourceLogger ContextResource = "Logger"
)

// ContextOperation is an enum for context operation
type ContextOperation string

const (
	// ContextOperationGet is for errors retrieving from context
	ContextOperationGet ContextOperation = "Get"
	// ContextOperationSet is for errors setting the context
	ContextOperationSet ContextOperation = "Set"
)

// ContextError is an error for working with context
type ContextError struct {
	Err       error
	Resource  ContextResource
	Operation ContextOperation
}

// Error provides the error as a string
func (e *ContextError) Error() string {
	return fmt.Sprintf(
		"could not %s %s on context with err: %s",
		e.Operation,
		e.Resource,
		e.Err,
	)
}

// Unwrap provides the underlying error
func (e *ContextError) Unwrap() error {
	return e.Err
}
