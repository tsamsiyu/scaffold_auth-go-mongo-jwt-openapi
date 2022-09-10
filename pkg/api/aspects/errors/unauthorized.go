package errors

import "encoding/json"

type UnauthorizedError struct {
	msg string
}

func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		msg: msg,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.msg
}

func (e *UnauthorizedError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": e.msg,
	})
}
