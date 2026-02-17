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

func TestProject_CRUD(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	setupAdminWithSession(t, env)

	var projectID string

	t.Run("Create", func(t *testing.T) {
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
		assert.False(t, result["is_archived"].(bool))

		t.Logf("Created project: %s", projectID)
	})

	t.Run("Read", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/projects/%s", projectID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, projectID, result["id"])
		assert.Equal(t, "Test Project", result["name"])
	})

	t.Run("Update", func(t *testing.T) {
		input := map[string]interface{}{
			"name":        "Updated Project Name",
			"description": "Updated description",
		}

		resp, err := env.Client.Put(fmt.Sprintf("/projects/%s", projectID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, "Updated Project Name", result["name"])
		assert.Equal(t, "Updated description", result["description"])

		t.Log("Project updated")
	})

	t.Run("List", func(t *testing.T) {
		resp, err := env.Client.Get("/projects")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})
}

func TestTask_Lifecycle(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	adminID := setupAdminWithSession(t, env)
	projectID, boardID := createProjectAndGetBoard(t, env)

	var taskID string

	t.Run("Create", func(t *testing.T) {
		input := map[string]interface{}{
			"title":       "Test Task",
			"description": "Testing task lifecycle",
			"status":      "backlog",
			"reporter_id": adminID,
		}

		resp, err := env.Client.Post(fmt.Sprintf("/boards/%s/tasks", boardID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		taskID = result["id"].(string)

		assert.Equal(t, "backlog", result["status"])
		assert.Equal(t, float64(0), result["position"])

		t.Logf("Task created: %s", taskID)
	})

	t.Run("Read", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/tasks/%s", taskID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, taskID, result["id"])
		assert.Equal(t, "Test Task", result["title"])
	})

	t.Run("Update", func(t *testing.T) {
		input := map[string]interface{}{
			"title":       "Updated Task Title",
			"description": "Updated description",
		}

		resp, err := env.Client.Put(fmt.Sprintf("/tasks/%s", taskID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, "Updated Task Title", result["title"])
	})

	t.Run("ListByBoard", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/boards/%s/tasks", boardID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})
}

func TestEpic_CRUD(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	setupAdminWithSession(t, env)
	projectID, _ := createProjectAndGetBoard(t, env)

	var epicID string

	t.Run("Create", func(t *testing.T) {
		input := fixtures.E2ENewEpicInput()

		resp, err := env.Client.Post(fmt.Sprintf("/projects/%s/epics", projectID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		epicID = result["id"].(string)

		assert.Equal(t, "Test Epic", result["name"])
		assert.Equal(t, projectID, result["project_id"])

		t.Logf("Epic created: %s", epicID)
	})

	t.Run("List", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/projects/%s/epics", projectID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("Read", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/epics/%s", epicID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, epicID, result["id"])
		assert.Equal(t, "Test Epic", result["name"])
	})

	t.Run("Update", func(t *testing.T) {
		input := map[string]interface{}{
			"name":        "Updated Epic",
			"description": "Updated description",
		}

		resp, err := env.Client.Put(fmt.Sprintf("/epics/%s", epicID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, "Updated Epic", result["name"])
	})
}

func TestBoard_CRUD(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	setupAdminWithSession(t, env)
	projectID, _ := createProjectAndGetBoard(t, env)

	t.Run("ListBoards", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1, "Should have at least one board")
	})

	var boardID string

	t.Run("CreateBoard", func(t *testing.T) {
		input := fixtures.E2ENewBoardInput()

		resp, err := env.Client.Post(fmt.Sprintf("/projects/%s/boards", projectID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		boardID = result["id"].(string)

		assert.Equal(t, "Test Board", result["name"])
	})

	t.Run("ReadBoard", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/boards/%s", boardID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		assert.Equal(t, boardID, result["id"])
		assert.Equal(t, "Test Board", result["name"])
	})
}

func TestLabel_CRUD(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	setupAdminWithSession(t, env)
	projectID, _ := createProjectAndGetBoard(t, env)

	var labelID string

	t.Run("Create", func(t *testing.T) {
		input := fixtures.E2ENewLabelInput()

		resp, err := env.Client.Post(fmt.Sprintf("/projects/%s/labels", projectID), input)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		labelID = result["id"].(string)

		assert.Equal(t, "Bug", result["name"])
		assert.Equal(t, "#FF0000", result["color"])

		t.Logf("Label created: %s", labelID)
	})

	t.Run("List", func(t *testing.T) {
		resp, err := env.Client.Get(fmt.Sprintf("/projects/%s/labels", projectID))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		ParseResponse(t, resp, &result)

		items := result["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 1)
	})
}

func setupAdminWithSession(t *testing.T, env *TestEnv) string {
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

	return fixtures.MemberID.String()
}

func createProjectAndGetBoard(t *testing.T, env *TestEnv) (string, string) {
	projectInput := fixtures.E2ENewProjectInput()
	resp, err := env.Client.Post("/projects", projectInput)
	require.NoError(t, err)

	var project map[string]interface{}
	ParseResponse(t, resp, &project)
	resp.Body.Close()

	projectID := project["id"].(string)

	resp, err = env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
	require.NoError(t, err)

	var boardsResult map[string]interface{}
	ParseResponse(t, resp, &boardsResult)
	resp.Body.Close()

	items := boardsResult["items"].([]interface{})
	require.GreaterOrEqual(t, len(items), 1, "Should have at least one board")

	board := items[0].(map[string]interface{})
	boardID := board["id"].(string)

	return projectID, boardID
}
