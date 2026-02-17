//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_List_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/projects")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestProject_Create_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	input := map[string]interface{}{
		"name": "Test Project",
		"key":  "TEST",
	}

	resp, err := env.Client.Post("/projects", input)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestMember_List_Unauthorized(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	resp, err := env.Client.Get("/members")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
