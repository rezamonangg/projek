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

func TestAdmin_GetStats_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/admin/stats")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAdmin_GetSettings_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	ctx := context.Background()

	err := env.ExecDB(ctx, `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, fixtures.CommunityID, "Test Community", "test-community")
	require.NoError(t, err)

	err = env.ExecDB(ctx, `
		INSERT INTO community_settings (id, community_id, allow_member_registration, require_email_verification, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, fixtures.CommunityID, fixtures.CommunityID, true, false)
	require.NoError(t, err)

	resp, err := env.Client.Get("/admin/settings")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAdmin_UpdateSettings_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	input := fixtures.E2ENewUpdateSettingsInput(false, true)

	resp, err := env.Client.Put("/admin/settings", input)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
