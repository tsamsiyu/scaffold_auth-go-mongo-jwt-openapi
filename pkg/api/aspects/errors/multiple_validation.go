package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type invalidEntry struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type MultipleValidationErrors struct {
	entries []invalidEntry
}

func NewMultipleValidationErrors(err error) *MultipleValidationErrors {
	if ozzoErr, ok := err.(validation.Errors); ok {
		return &MultipleValidationErrors{
			entries: mapOzzoValidationErrors(ozzoErr),
		}
	}

	panic(fmt.Sprintf("Unknown error type to cast into MultipleValidationErrors: %T", err))
}

func NewMultipleValidationInputError(err error) *InputError {
	return NewInputError(NewMultipleValidationErrors(err))
}

func (e *MultipleValidationErrors) Error() string {
	var entriesStr []string
	for _, entry := range e.entries {
		entriesStr = append(entriesStr, fmt.Sprintf("%s ", entry.Message))
	}

	return strings.Join(entriesStr, "; ")
}

func (e MultipleValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"tag":     "validation",
		"entries": e.entries,
	})
}

func mapOzzoValidationErrors(src validation.Errors) []invalidEntry {
	entries := make([]invalidEntry, 0)

	flattenValidationErrorsRecursively("", &src, &entries)

	return entries
}

func flattenValidationErrorsRecursively(prefix string, src *validation.Errors, dest *[]invalidEntry) {
	for k, v := range *src {
		var chainKey string
		if prefix != "" {
			chainKey = prefix + "." + k
		} else {
			chainKey = k
		}

		if nested, ok := v.(validation.Errors); ok {
			flattenValidationErrorsRecursively(chainKey, &nested, dest)
			continue
		}

		if edge, ok := v.(error); ok {
			*dest = append(*dest, invalidEntry{
				Path:    chainKey,
				Message: edge.Error(),
			})
		}
	}
}
