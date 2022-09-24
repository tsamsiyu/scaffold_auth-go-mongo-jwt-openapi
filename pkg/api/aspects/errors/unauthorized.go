package errors

import (
	"encoding/json"
	"fmt"
)

type UnauthorizedError struct {
	reason string
}

func NewUnauthorizedError(reason string) *UnauthorizedError {
	return &UnauthorizedError{
		reason: reason,
	}
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("User unauthorized with reason %s", e.reason)
}

func (e *UnauthorizedError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": "User unauthorized",
		"reason":  e.reason,
	})
}
