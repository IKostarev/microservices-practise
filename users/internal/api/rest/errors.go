package rest

type ErrorCode string

const (
	ErrCodeInternalAPIError       ErrorCode = "INTERNAL_API_ERROR"
	ErrCodeRequestValidationError ErrorCode = "USER_REQUEST_VALIDATION_ERROR"
	ErrCodeBadRequest             ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound               ErrorCode = "NOT_FOUND"
	ErrCodeInvalidJsonFormat      ErrorCode = "JSON_FORMAT_ERROR"
)

type ApiError struct {
	Message string
	ErrCode ErrorCode
}

func NewApiError(message string, errCode ErrorCode) *ApiError {
	return &ApiError{
		Message: message,
		ErrCode: errCode,
	}
}

func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) As(target interface{}) bool {
	t, ok := target.(**ApiError)
	if !ok {
		return false
	}

	*t = e

	return true
}

var (
	ErrInternalApi                = NewApiError("an error occurred while processing the request", ErrCodeInternalAPIError)
	ErrBadRequest                 = NewApiError("bad request", ErrCodeBadRequest)
	ErrUsernameOrEmailAlreadyUsed = NewApiError("username or email already used", ErrCodeBadRequest)
	ErrWrongCredentials           = NewApiError("wrong credentials", ErrCodeBadRequest)
	ErrNotFound                   = NewApiError("not found", ErrCodeNotFound)
)

func (e *ApiError) IsRequestValidationError() bool {
	return e.ErrCode == ErrCodeRequestValidationError
}
