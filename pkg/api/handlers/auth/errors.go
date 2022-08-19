package auth

import (
	"apart-deal-api/pkg/store/user"

	apiErr "apart-deal-api/pkg/api/aspects/errors"
	authDomain "apart-deal-api/pkg/domain/auth"
)

func mapError(err error) error {
	if _, ok := err.(*user.UserDuplicateError); ok {
		return apiErr.NewConflictError("Such user already exists")
	}

	if _, ok := err.(*authDomain.UserNotFound); ok {
		return apiErr.NewNotFoundError("User not found")
	}

	return err
}
