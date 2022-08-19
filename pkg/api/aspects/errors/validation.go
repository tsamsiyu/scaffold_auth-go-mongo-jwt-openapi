package errors

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ValidationError struct {
	error
}

func NewValidationError(err error) *ValidationError {
	if vErr, ok := err.(validation.Errors); ok {
		return &ValidationError{NewMultipleValidationErrors(vErr)}
	}

	return &ValidationError{error: err}
}

func (e *ValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.error)
}
