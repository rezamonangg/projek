package common

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("conflict")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func ErrToCode(err error) string {
	switch {
	case errors.Is(err, ErrNotFound):
		return "NOT_FOUND"
	case errors.Is(err, ErrUnauthorized):
		return "UNAUTHORIZED"
	case errors.Is(err, ErrForbidden):
		return "FORBIDDEN"
	case errors.Is(err, ErrBadRequest):
		return "BAD_REQUEST"
	case errors.Is(err, ErrConflict):
		return "CONFLICT"
	default:
		return "INTERNAL_ERROR"
	}
}
