package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := &AppError{
			Code:    "TEST_CODE",
			Message: "test message",
			Err:     assert.AnError,
		}
		assert.Equal(t, "test message", err.Error())
	})

	t.Run("unwrap", func(t *testing.T) {
		err := &AppError{
			Code:    "TEST_CODE",
			Message: "test message",
			Err:     assert.AnError,
		}
		assert.True(t, errors.Is(err, assert.AnError))
	})
}

func TestNewError(t *testing.T) {
	t.Run("creates error with all fields", func(t *testing.T) {
		err := NewError("CODE", "message", assert.AnError)
		assert.Equal(t, "CODE", err.Code)
		assert.Equal(t, "message", err.Message)
		assert.Equal(t, assert.AnError, err.Err)
	})

	t.Run("creates error without underlying error", func(t *testing.T) {
		err := NewError("CODE", "message", nil)
		assert.Equal(t, "CODE", err.Code)
		assert.Equal(t, "message", err.Message)
		assert.Nil(t, err.Err)
	})
}

func TestErrToCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"not found", ErrNotFound, "NOT_FOUND"},
		{"unauthorized", ErrUnauthorized, "UNAUTHORIZED"},
		{"forbidden", ErrForbidden, "FORBIDDEN"},
		{"bad request", ErrBadRequest, "BAD_REQUEST"},
		{"conflict", ErrConflict, "CONFLICT"},
		{"unknown error", assert.AnError, "INTERNAL_ERROR"},
		{"wrapped not found", errors.New("wrapped: not found"), "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrToCode(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	assert.Equal(t, "resource not found", ErrNotFound.Error())
	assert.Equal(t, "unauthorized", ErrUnauthorized.Error())
	assert.Equal(t, "forbidden", ErrForbidden.Error())
	assert.Equal(t, "bad request", ErrBadRequest.Error())
	assert.Equal(t, "conflict", ErrConflict.Error())
	assert.Equal(t, "internal server error", ErrInternalServer.Error())
	assert.Equal(t, "invalid credentials", ErrInvalidCredentials.Error())
	assert.Equal(t, "user already exists", ErrUserExists.Error())
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "token expired", ErrExpiredToken.Error())
}
