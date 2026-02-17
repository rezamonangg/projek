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

	passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3/ItB/XBG/eCknfIrqS6"
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

	var result map[string]interface{}
	ParseResponse(t, resp, &result)

	assert.NotEmpty(t, result["session_id"], "Should return session ID")

	member := result["member"].(map[string]interface{})
	assert.Equal(t, "admin@test.com", member["email"])
	assert.Equal(t, "admin", member["role"])

	sessionCookie := env.Client.GetSessionCookie()
	assert.NotNil(t, sessionCookie, "Session cookie should be set")
	assert.NotEmpty(t, sessionCookie.Value, "Session cookie should have value")

	t.Logf("Login successful, session: %s", result["session_id"])
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

	passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3/ItB/XBG/eCknfIrqS6"
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

func TestAuth_Logout_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	registerAndLogin(t, env)

	sessionCookie := env.Client.GetSessionCookie()
	require.NotNil(t, sessionCookie)

	resp, err := env.Client.Post("/auth/logout", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ctx := context.Background()
	exists := env.Redis.Exists(ctx, "session:"+sessionCookie.Value).Val()
	assert.Equal(t, int64(0), exists, "Session should be removed from Redis")
}

func TestAuth_Logout_WithoutSession(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Post("/auth/logout", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_ProtectedEndpoint_WithoutAuth(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ProtectedEndpoint_WithAuth(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	registerAndLogin(t, env)

	resp, err := env.Client.Get("/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	ParseResponse(t, resp, &result)

	assert.Equal(t, "admin@test.com", result["email"])
	assert.Equal(t, "admin", result["role"])
}

func TestAuth_ProtectedEndpoint_AfterLogout(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	registerAndLogin(t, env)

	resp, _ := env.Client.Post("/auth/logout", nil)
	resp.Body.Close()

	resp, err := env.Client.Get("/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ProtectedEndpoint_ExpiredSession(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	registerAndLogin(t, env)

	sessionCookie := env.Client.GetSessionCookie()
	require.NotNil(t, sessionCookie)

	ctx := context.Background()
	env.Redis.Del(ctx, "session:"+sessionCookie.Value)

	resp, err := env.Client.Get("/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_CompleteFlow(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	t.Log("Step 1: Create community and member in database")
	err := env.ExecDB(ctx, `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, fixtures.CommunityID, "Test Community", "test-community")
	require.NoError(t, err)

	passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3/ItB/XBG/eCknfIrqS6"
	err = env.ExecDB(ctx, `
		INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, fixtures.MemberID, fixtures.CommunityID, "admin@test.com", passwordHash, "Admin", "User", "admin", true)
	require.NoError(t, err)

	t.Log("Step 2: Login as admin")
	loginInput := fixtures.E2ENewLoginInput("admin@test.com", "password123")

	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	t.Log("Step 3: Access protected endpoint")
	resp, err = env.Client.Get("/auth/me")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var meResult map[string]interface{}
	ParseResponse(t, resp, &meResult)

	assert.Equal(t, "admin@test.com", meResult["email"])
	assert.Equal(t, "admin", meResult["role"])
	resp.Body.Close()

	t.Log("Step 4: Logout")
	resp, err = env.Client.Post("/auth/logout", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	t.Log("Step 5: Verify protected endpoint is now inaccessible")
	resp, err = env.Client.Get("/auth/me")
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	t.Log("Complete auth flow verified")
}

func registerAndLogin(t *testing.T, env *TestEnv) {
	ctx := context.Background()

	err := env.ExecDB(ctx, `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, fixtures.CommunityID, "Test Community", "test-community")
	require.NoError(t, err)

	passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3/ItB/XBG/eCknfIrqS6"
	err = env.ExecDB(ctx, `
		INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, fixtures.MemberID, fixtures.CommunityID, "admin@test.com", passwordHash, "Admin", "User", "admin", true)
	require.NoError(t, err)

	loginInput := fixtures.E2ENewLoginInput("admin@test.com", "password123")
	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	resp.Body.Close()

	sessionCookie := env.Client.GetSessionCookie()
	require.NotNil(t, sessionCookie, "Should have session cookie after login")
}
