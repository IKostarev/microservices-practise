package app_errors

import "errors"

var (
	ErrUsernameOrEmailIsUsed           = errors.New("username or email already used")
	ErrNotFound                        = errors.New("not found")
	ErrWrongCredentials                = errors.New("wrong credentials")
	ErrIncorrectOldPassword            = errors.New("incorrect old password")
	ErrPassAndConfirmationDoesNotMatch = errors.New("password and confirmation does not match")
)
