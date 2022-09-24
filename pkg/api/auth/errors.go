package auth

type TokenInvalidError struct {
}

func (e *TokenInvalidError) Error() string {
	return "Token is invalid"
}

type TokenDoesNotExistError struct {
}

func (e *TokenDoesNotExistError) Error() string {
	return "Token does not exist"
}

type TokenExpiredError struct {
}

func (e *TokenExpiredError) Error() string {
	return "Token has expired"
}

type UserNotConfirmedError struct {
	error
}

func (e *UserNotConfirmedError) Error() string {
	return "User is not confirmed"
}

type InvalidPasswordError struct {
	error
}

func (e *InvalidPasswordError) Error() string {
	return "Invalid password"
}

type NoSuchUserError struct {
	error
}

func (e *NoSuchUserError) Error() string {
	return "No such user"
}
