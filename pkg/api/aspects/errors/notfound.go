package errors

import "encoding/json"

type NotFoundError struct {
	msg string
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{msg}
}

func (e *NotFoundError) Error() string {
	return e.msg
}

func (e *NotFoundError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": e.msg,
	})
}
