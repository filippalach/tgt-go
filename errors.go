package tgtg

import (
	"fmt"
	"net/http"

	"github.com/logrusorgru/aurora"
)

// ArgumentError is an error that represents an issue with an input to Too Good To Go
// API client. Identifies input parameter and the cause.
type ArgumentError struct {
	// argument specifies function's argument name.
	argument string

	// reason specifies reason of the error.
	reason string
}

var _ error = &ArgumentError{}

// Error implements error interface's method.
func (e *ArgumentError) Error() string {
	return fmt.Sprintf("%s argument is invalid (reason: %s)", e.argument, e.reason)
}

// NewArgumentError creates an ArgumentError.
func NewArgumentError(argument, reason string) *ArgumentError {
	return &ArgumentError{
		argument: argument,
		reason:   reason,
	}
}

// ErrorResponse represents error response returned from Too Good To Go API.
type ErrorResponse struct {
	// HTTP response that caused this error.
	Response *http.Response

	// Error body from Too Good To Go API.
	Errors []Error `json:"errors"`
}

var _ error = &ErrorResponse{}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(`
	method: %v
	url:    %v
	code:   %d
	errors: %v`,
		aurora.Green(r.Response.Request.Method),
		aurora.Yellow(r.Response.Request.URL),
		aurora.Red(r.Response.StatusCode),
		aurora.Red(r.Errors))
}

// An Error represent the error returned from Too Good to Go API.
type Error struct {
	// Code field is always set by Too Good To Go API.
	Code string `json:"code"`

	// Message is optional and not always returned.
	Message string `json:"message"`
}
