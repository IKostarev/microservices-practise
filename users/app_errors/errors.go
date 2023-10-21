package app_errors

import "errors"

var (
	ErrUsernameOrEmailIsUsed           = errors.New("username or email already used")
	ErrWrongPassword                   = errors.New("wrong password")
	ErrIncorrectOldPassword            = errors.New("incorrect old password")
	ErrPassAndConfirmationDoesNotMatch = errors.New("password and confirmation does not match")
)
