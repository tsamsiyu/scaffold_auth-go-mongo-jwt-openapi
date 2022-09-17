package auth

import (
	"apart-deal-api/pkg/api/auth"
	"apart-deal-api/pkg/store/user"

	apiErr "apart-deal-api/pkg/api/aspects/errors"
	authDomain "apart-deal-api/pkg/domain/auth"
)

func mapError(err error) error {
	if _, ok := err.(*user.UserDuplicateError); ok {
		return apiErr.NewConflictError("Such user already exists")
	}

	if _, ok := err.(*auth.InvalidPasswordError); ok {
		return apiErr.NewSimpleValidationInputError("Password is invalid", "invalid_pass")
	}

	if _, ok := err.(*authDomain.UserNotFound); ok {
		return apiErr.NewNotFoundError("User not found")
	}

	if _, ok := err.(*authDomain.CouldNotConfirmError); ok {
		return apiErr.NewInputError(
			apiErr.NewSimpleValidationError("Could not confirm this user", "unconfirmable"),
		)
	}

	if _, ok := err.(*authDomain.ConfirmationCodeMismatchError); ok {
		return apiErr.NewInputError(
			apiErr.NewSimpleValidationError("Confirmation code mismatched", "code_mismatch"),
		)
	}

	if _, ok := err.(*auth.UserNotConfirmedError); ok {
		return apiErr.NewUnauthorizedError(err.Error())
	}

	return err
}
