package auth

type UserNotFound struct {
	error
}

type CouldNotConfirmError struct {
	error
}
