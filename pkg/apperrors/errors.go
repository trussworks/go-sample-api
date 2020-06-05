package apperrors

import (
	"encoding/json"
	"fmt"

	"bin/bork/pkg/models"
)

// QueryOperation provides a set of operations that can fail
type QueryOperation string

const (
	// QueryCreate is for failures when creating a resource
	QueryCreate QueryOperation = "Create"
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

// Validations maps attributes to validation messages
type Validations map[string]string

// Map directly returns a map in case implementation of Validations changes
func (v Validations) Map() map[string]string {
	return v
}

// NewValidationError returns a validation error with fields instantiated
func NewValidationError(err error, resource interface{}, resourceID string) ValidationError {
	return ValidationError{
		Err:         err,
		Validations: Validations{},
		Resource:    resource,
		ResourceID:  resourceID,
	}
}

// ValidationError is a typed error for issues with validation
type ValidationError struct {
	Err         error
	Validations Validations
	Resource    interface{}
	ResourceID  string
}

// WithValidation allows a failed validation message be added to the ValidationError
func (e ValidationError) WithValidation(key string, message string) {
	e.Validations[key] = message
}

// Error provides the error as a string
func (e *ValidationError) Error() string {
	data, err := json.Marshal(e.Validations)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Could not validate %T %s: %s", e.Resource, e.ResourceID, string(data))
}

// Unwrap provides the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// MethodNotAllowedError is a typed error for HTTP methods not allowed
type MethodNotAllowedError struct {
	Method string
}

// Error provides the error as a string
func (e *MethodNotAllowedError) Error() string {
	return fmt.Sprintf(
		"Method %s not allowed",
		e.Method,
	)
}

// UnknownRouteError is an error for unknown routes
type UnknownRouteError struct {
	Path string
}

// Error provides the error as a string
func (e *UnknownRouteError) Error() string {
	return fmt.Sprintf(
		"Route %s unknown",
		e.Path,
	)
}

// BadRequestError is a typed error for bad request content
type BadRequestError struct {
	Err error
}

// Error provides the error as a string
func (e *BadRequestError) Error() string {
	return fmt.Sprintf(
		"Request could not understood: %v",
		e.Err,
	)
}

// Unwrap provides the underlying error
func (e *BadRequestError) Unwrap() error {
	return e.Err
}
