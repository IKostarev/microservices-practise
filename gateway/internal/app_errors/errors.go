package app_errors

import (
	"errors"
	"fmt"
)

var (
	ErrUsernameOrEmailIsUsed           = errors.New("username or email already used")
	ErrNotFound                        = errors.New("not found")
	ErrWrongCredentials                = errors.New("wrong credentials")
	ErrIncorrectOldPassword            = errors.New("incorrect old password")
	ErrPassAndConfirmationDoesNotMatch = errors.New("password and confirmation does not match")
	ErrNoUserInContext                 = errors.New("no user in context")
)

type UserIDMismatchError struct {
	Operation string
	ContextID int
	UserDTOID int
}

// Error implements the error interface for UserIDMismatchError
func (e *UserIDMismatchError) Error() string {
	return fmt.Sprintf("%s context ID [%d], DTO ID [%d]", e.Operation, e.ContextID, e.UserDTOID)
}

func NewUserIDMismatchError(
	Operation string,
	ContextID int,
	UserDTOID int,
) *UserIDMismatchError {
	return &UserIDMismatchError{
		Operation: Operation,
		ContextID: ContextID,
		UserDTOID: UserDTOID,
	}
}
