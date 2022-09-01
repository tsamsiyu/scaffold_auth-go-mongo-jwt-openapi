package errors

import "encoding/json"

type SimpleValidationError struct {
	msg string
	tag string
}

func NewSimpleValidationError(msg string, tag string) *SimpleValidationError {
	return &SimpleValidationError{
		msg: msg,
		tag: tag,
	}
}

func (e *SimpleValidationError) Error() string {
	return e.msg
}

func (e *SimpleValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"message": e.msg,
		"tag":     e.tag,
	})
}
