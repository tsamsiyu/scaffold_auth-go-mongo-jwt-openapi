package auth

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
