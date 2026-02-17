//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/monachy/projek/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_Login_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	err := env.ExecDB(ctx, `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, fixtures.CommunityID, "Test Community", "test-community")
	require.NoError(t, err)

	passwordHash := "$2a$10$sKr8GxHBJAbOdab9Ma.6NOVV9XFORV6bKg1VOpH1guE7rlv4SucO."
	err = env.ExecDB(ctx, `
		INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, fixtures.MemberID, fixtures.CommunityID, "admin@test.com", passwordHash, "Admin", "User", "admin", true)
	require.NoError(t, err)

	loginInput := fixtures.E2ENewLoginInput("admin@test.com", "password123")

	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := ParseDataResponse(t, resp)

	assert.NotEmpty(t, data["session_id"], "Should return session ID")

	member := data["member"].(map[string]interface{})
	assert.Equal(t, "admin@test.com", member["email"])
	assert.Equal(t, "admin", member["role"])

	t.Logf("Login successful, session: %s", data["session_id"])
}

func TestAuth_Login_WrongPassword(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	err := env.ExecDB(ctx, `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, fixtures.CommunityID, "Test Community", "test-community")
	require.NoError(t, err)

	passwordHash := "$2a$10$sKr8GxHBJAbOdab9Ma.6NOVV9XFORV6bKg1VOpH1guE7rlv4SucO."
	err = env.ExecDB(ctx, `
		INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, fixtures.MemberID, fixtures.CommunityID, "admin@test.com", passwordHash, "Admin", "User", "admin", true)
	require.NoError(t, err)

	loginInput := fixtures.E2ENewLoginInput("admin@test.com", "wrongpassword")

	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_Login_NonExistentUser(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	loginInput := fixtures.E2ENewLoginInput("nobody@test.com", "password123")

	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_Login_ValidationErrors(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	tests := []struct {
		name       string
		input      fixtures.E2ELoginInput
		expectCode int
	}{
		{
			name:       "missing email",
			input:      fixtures.E2ENewLoginInput("", "password"),
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			input:      fixtures.E2ENewLoginInput("admin@test.com", ""),
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := env.Client.Post("/auth/login", tt.input)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectCode, resp.StatusCode)
		})
	}
}

func TestAuth_ProtectedEndpoint_WithoutAuth(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHealth_Endpoints(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("Health", func(t *testing.T) {
		resp, err := env.Client.Get("/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Ready", func(t *testing.T) {
		resp, err := env.Client.Get("/ready")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
