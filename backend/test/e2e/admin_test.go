//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/monachy/projek/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdmin_GetStats_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	setupAdminWithData(t, env)

	resp, err := env.Client.Get("/admin/stats")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	ParseResponse(t, resp, &result)

	assert.Contains(t, result, "total_members")
	assert.Contains(t, result, "total_projects")
	assert.Contains(t, result, "total_tasks")

	t.Logf("Stats: members=%v, projects=%v, tasks=%v",
		result["total_members"], result["total_projects"], result["total_tasks"])
}

func TestAdmin_GetStats_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/admin/stats")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAdmin_GetSettings_Success(t *testing.T) {
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

	err = env.ExecDB(ctx, `
		INSERT INTO community_settings (id, community_id, allow_member_registration, require_email_verification, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, fixtures.CommunityID, fixtures.CommunityID, true, false)
	require.NoError(t, err)

	loginAndSetSession(t, env)

	resp, err := env.Client.Get("/admin/settings")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	ParseResponse(t, resp, &result)

	assert.Equal(t, true, result["allow_member_registration"])
	assert.Equal(t, false, result["require_email_verification"])
}

func TestAdmin_UpdateSettings_Success(t *testing.T) {
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

	err = env.ExecDB(ctx, `
		INSERT INTO community_settings (id, community_id, allow_member_registration, require_email_verification, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, fixtures.CommunityID, fixtures.CommunityID, true, false)
	require.NoError(t, err)

	loginAndSetSession(t, env)

	input := fixtures.E2ENewUpdateSettingsInput(false, true)

	resp, err := env.Client.Put("/admin/settings", input)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	ParseResponse(t, resp, &result)

	assert.Equal(t, false, result["allow_member_registration"])
	assert.Equal(t, true, result["require_email_verification"])
}

func TestAdmin_UpdateSettings_ForbiddenForNonAdmin(t *testing.T) {
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
	`, fixtures.MemberID, fixtures.CommunityID, "member@test.com", passwordHash, "Regular", "Member", "member", true)
	require.NoError(t, err)

	err = env.ExecDB(ctx, `
		INSERT INTO community_settings (id, community_id, allow_member_registration, require_email_verification, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, fixtures.CommunityID, fixtures.CommunityID, true, false)
	require.NoError(t, err)

	loginInput := fixtures.E2ENewLoginInput("member@test.com", "password123")
	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	resp.Body.Close()

	input := fixtures.E2ENewUpdateSettingsInput(false, true)

	resp, err = env.Client.Put("/admin/settings", input)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestAdmin_CompleteWorkflow(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	var projectID, boardID string

	t.Run("Phase1_Setup", func(t *testing.T) {
		setupAdminWithData(t, env)
		t.Log("Admin setup complete")
	})

	t.Run("Phase2_CreateProject", func(t *testing.T) {
		input := fixtures.E2ENewProjectInput()

		resp, err := env.Client.Post("/projects", input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		projectID = result["id"].(string)

		assert.Equal(t, "Test Project", result["name"])
		assert.Equal(t, "TEST", result["key"])

		t.Logf("Project created: %s", projectID)
	})

	t.Run("Phase3_VerifyDefaultBoard", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		require.GreaterOrEqual(t, len(items), 1, "Should have at least 1 board")

		board := items[0].(map[string]interface{})
		boardID = board["id"].(string)

		t.Logf("Default board: %s", boardID)
	})

	t.Run("Phase4_GetDashboardStats", func(t *testing.T) {
		resp, err := env.Client.Get("/admin/stats")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		totalProjects := result["total_projects"].(float64)
		assert.GreaterOrEqual(t, totalProjects, float64(1), "At least 1 project")

		t.Logf("Dashboard stats: %+v", result)
	})

	t.Log("Complete admin workflow verified")
}

func setupAdminWithData(t *testing.T, env *TestEnv) {
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

	loginAndSetSession(t, env)
}

func loginAndSetSession(t *testing.T, env *TestEnv) {
	loginInput := fixtures.E2ENewLoginInput("admin@test.com", "password123")
	resp, err := env.Client.Post("/auth/login", loginInput)
	require.NoError(t, err)
	resp.Body.Close()

	sessionCookie := env.Client.GetSessionCookie()
	require.NotNil(t, sessionCookie, "Should have session cookie after login")
}
