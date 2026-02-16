package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"short password", "pass", false},
		{"empty password", "", false},
		{"long password", "verylongpassword1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"correct password", "password123", hash, true},
		{"wrong password", "wrongpassword", hash, false},
		{"empty password", "", hash, false},
		{"invalid hash", "password123", "invalidhash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashPasswordUnique(t *testing.T) {
	password := "password123"
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	hash2, err := HashPassword(password)
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "same password should produce different hashes due to salt")
	assert.True(t, CheckPassword(password, hash1))
	assert.True(t, CheckPassword(password, hash2))
}
