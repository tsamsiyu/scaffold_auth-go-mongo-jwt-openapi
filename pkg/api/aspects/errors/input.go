package errors

import (
	"encoding/json"
)

type InputError struct {
	error
}

func NewInputError(err error) *InputError {
	return &InputError{error: err}
}

func (e *InputError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.error)
}
