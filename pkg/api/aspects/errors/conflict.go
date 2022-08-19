package errors

import "encoding/json"

type ConflictError struct {
	msg string
}

func NewConflictError(msg string) *ConflictError {
	return &ConflictError{msg}
}

func (e *ConflictError) Error() string {
	return e.msg
}

func (e *ConflictError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": e.msg,
	})
}
