package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValidStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,min=18"`
}

type MemberInput struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
	FirstName string `validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			input: ValidStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			input: ValidStruct{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: ValidStruct{
				Name:  "John",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "age below minimum",
			input: ValidStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   16,
			},
			wantErr: true,
		},
		{
			name: "member create input valid",
			input: MemberInput{
				Email:     "john@example.com",
				Password:  "password123",
				FirstName: "John",
			},
			wantErr: false,
		},
		{
			name: "member create input invalid email",
			input: MemberInput{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
			},
			wantErr: true,
		},
		{
			name: "member create input short password",
			input: MemberInput{
				Email:     "john@example.com",
				Password:  "short",
				FirstName: "John",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	input := ValidStruct{
		Name:  "",
		Email: "invalid",
		Age:   16,
	}

	err := ValidateStruct(input)
	require.Error(t, err)

	errors := ValidationError(err)
	assert.NotEmpty(t, errors)
	assert.Contains(t, errors, "Name")
	assert.Contains(t, errors, "Email")
	assert.Contains(t, errors, "Age")
}
