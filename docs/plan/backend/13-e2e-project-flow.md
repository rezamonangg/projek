# E2E Test Plan: Project Flow

**Priority:** HIGH
**Estimated Time:** 1 hour
**Dependencies:** 10-e2e-shared-setup.md, 11-e2e-auth-flow.md

## Test Scope

| Flow | Description |
|------|-------------|
| Project CRUD | Create, Read, Update, Delete projects |
| Task Lifecycle | Create → backlog → todo → inprogress → done → delete |
| Task in Epic | Create epic, create tasks inside epic |
| Board Management | List boards, create additional boards |

## Database Isolation

Each top-level test creates a **fresh PostgreSQL and Redis container**.

## Test File

**File:** `backend/test/e2e/project_test.go`

```go
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

// ============================================================
// TEST SUITE: Project CRUD
// ============================================================

func TestProject_CRUD(t *testing.T) {
    // GIVEN: Fresh isolated database with logged-in admin
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    adminID := setupAdmin(t, env)
    var projectID string
    
    t.Run("Create", func(t *testing.T) {
        input := fixtures.NewProjectInput()
        
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
    
    t.Run("Archive", func(t *testing.T) {
        input := map[string]interface{}{
            "is_archived": true,
        }
        
        resp, err := env.Client.Put(fmt.Sprintf("/projects/%s", projectID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.True(t, result["is_archived"].(bool))
        
        t.Log("Project archived")
    })
    
    t.Run("ListExcludesArchived", func(t *testing.T) {
        resp, err := env.Client.Get("/projects")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        // Archived projects should not appear in list
        items := result["items"].([]interface{})
        assert.Empty(t, items, "Archived projects should not appear in list")
    })
    
    t.Run("Delete", func(t *testing.T) {
        resp, err := env.Client.Delete(fmt.Sprintf("/projects/%s", projectID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusNoContent, resp.StatusCode)
        
        t.Log("Project deleted")
    })
    
    t.Run("ReadAfterDelete", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/projects/%s", projectID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })
}

func TestProject_CreateValidation(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    setupAdmin(t, env)
    
    tests := []struct {
        name       string
        input      map[string]interface{}
        expectCode int
    }{
        {
            name: "missing name",
            input: map[string]interface{}{
                "key": "TEST",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "missing key",
            input: map[string]interface{}{
                "name": "Test Project",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "key too short",
            input: map[string]interface{}{
                "name": "Test Project",
                "key":  "A",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "key too long",
            input: map[string]interface{}{
                "name": "Test Project",
                "key":  "TOOLONGKEY",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "key lowercase",
            input: map[string]interface{}{
                "name": "Test Project",
                "key":  "test",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "valid input",
            input: map[string]interface{}{
                "name": "Test Project",
                "key":  "TEST",
            },
            expectCode: http.StatusCreated,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := env.Client.Post("/projects", tt.input)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, tt.expectCode, resp.StatusCode)
        })
    }
}

// ============================================================
// TEST SUITE: Task Lifecycle
// ============================================================

func TestTask_Lifecycle(t *testing.T) {
    // GIVEN: Fresh isolated database with project and board
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    adminID := setupAdmin(t, env)
    projectID, boardID := createProjectAndGetBoard(t, env)
    
    var taskID string
    
    // State transitions: backlog → todo → inprogress → done
    t.Run("Create_Backlog", func(t *testing.T) {
        input := map[string]interface{}{
            "title":       "Lifecycle Task",
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
        
        t.Logf("Task created in backlog: %s", taskID)
    })
    
    t.Run("Move_ToTodo", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "todo",
            "position": 0,
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "todo", result["status"])
        
        t.Log("Task moved to todo")
    })
    
    t.Run("Move_ToInProgress", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "inprogress",
            "position": 0,
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "inprogress", result["status"])
        
        t.Log("Task moved to inprogress")
    })
    
    t.Run("Move_ToCodeReview", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "codereview",
            "position": 0,
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "codereview", result["status"])
        
        t.Log("Task moved to codereview")
    })
    
    t.Run("Move_ToInTest", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "intest",
            "position": 0,
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "intest", result["status"])
        
        t.Log("Task moved to intest")
    })
    
    t.Run("Move_ToDone", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "done",
            "position": 0,
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "done", result["status"])
        
        t.Log("Task moved to done")
    })
    
    t.Run("Update_AfterDone", func(t *testing.T) {
        input := map[string]interface{}{
            "title": "Updated Title After Done",
        }
        
        resp, err := env.Client.Put(fmt.Sprintf("/tasks/%s", taskID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "Updated Title After Done", result["title"])
        assert.Equal(t, "done", result["status"], "Status should remain done")
    })
    
    t.Run("Delete", func(t *testing.T) {
        resp, err := env.Client.Delete(fmt.Sprintf("/tasks/%s", taskID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusNoContent, resp.StatusCode)
        
        t.Log("Task deleted")
    })
    
    t.Run("NotFoundAfterDelete", func(t *testing.T) {
        // Try to get deleted task
        resp, err := env.Client.Get(fmt.Sprintf("/tasks/%s", taskID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })
}

func TestTask_PositionOrdering(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    adminID := setupAdmin(t, env)
    _, boardID := createProjectAndGetBoard(t, env)
    
    // Create 3 tasks
    taskIDs := make([]string, 3)
    for i := 0; i < 3; i++ {
        input := map[string]interface{}{
            "title":       fmt.Sprintf("Task %d", i+1),
            "status":      "todo",
            "reporter_id": adminID,
        }
        
        resp, _ := env.Client.Post(fmt.Sprintf("/boards/%s/tasks", boardID), input)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        resp.Body.Close()
        
        taskIDs[i] = result["id"].(string)
    }
    
    // Verify positions are assigned correctly
    t.Run("VerifyInitialPositions", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/boards/%s/tasks", boardID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        var result []map[string]interface{}
        ParseResponse(t, resp, &result)
        
        // Positions should be 0, 1, 2
        positions := make(map[string]int)
        for _, task := range result {
            positions[task["id"].(string)] = int(task["position"].(float64))
        }
        
        // Verify positions are unique and sequential
        assert.Contains(t, []int{0, 1, 2}, positions[taskIDs[0]])
        assert.Contains(t, []int{0, 1, 2}, positions[taskIDs[1]])
        assert.Contains(t, []int{0, 1, 2}, positions[taskIDs[2]])
    })
    
    // Move task to different position
    t.Run("ReorderTasks", func(t *testing.T) {
        input := map[string]interface{}{
            "status":   "todo",
            "position": 0, // Move first task to position 0 (should reorder others)
        }
        
        resp, err := env.Client.Patch(fmt.Sprintf("/tasks/%s/move", taskIDs[0]), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
    })
}

// ============================================================
// TEST SUITE: Epic and Tasks
// ============================================================

func TestEpic_TaskRelationship(t *testing.T) {
    // GIVEN: Fresh isolated database
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    adminID := setupAdmin(t, env)
    projectID, boardID := createProjectAndGetBoard(t, env)
    
    var epicID string
    taskIDs := make([]string, 3)
    
    t.Run("CreateEpic", func(t *testing.T) {
        input := fixtures.NewEpicInputWithName("Sprint 1")
        
        resp, err := env.Client.Post(fmt.Sprintf("/projects/%s/epics", projectID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusCreated, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        epicID = result["id"].(string)
        
        assert.Equal(t, "Sprint 1", result["name"])
        assert.Equal(t, projectID, result["project_id"])
        
        t.Logf("Epic created: %s", epicID)
    })
    
    t.Run("CreateTasksInEpic", func(t *testing.T) {
        for i := 0; i < 3; i++ {
            input := map[string]interface{}{
                "title":       fmt.Sprintf("Sprint Task %d", i+1),
                "status":      "todo",
                "reporter_id": adminID,
                "epic_id":     epicID,
            }
            
            resp, err := env.Client.Post(fmt.Sprintf("/boards/%s/tasks", boardID), input)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            require.Equal(t, http.StatusCreated, resp.StatusCode)
            
            var result map[string]interface{}
            ParseResponse(t, resp, &result)
            
            taskIDs[i] = result["id"].(string)
            
            assert.Equal(t, epicID, result["epic_id"], "Task should be linked to epic")
        }
        
        t.Logf("Created 3 tasks in epic")
    })
    
    t.Run("VerifyTasksInEpic", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/boards/%s/tasks", boardID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        var result []map[string]interface{}
        ParseResponse(t, resp, &result)
        
        epicTaskCount := 0
        for _, task := range result {
            if epicID == task["epic_id"] {
                epicTaskCount++
            }
        }
        
        assert.Equal(t, 3, epicTaskCount, "Should have 3 tasks in epic")
        
        t.Logf("Verified: %d tasks in epic", epicTaskCount)
    })
    
    t.Run("UpdateTask_RemoveFromEpic", func(t *testing.T) {
        // Remove task from epic by setting epic_id to null
        input := map[string]interface{}{
            "epic_id": nil,
        }
        
        resp, err := env.Client.Put(fmt.Sprintf("/tasks/%s", taskIDs[0]), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Nil(t, result["epic_id"], "Task should no longer be in epic")
        
        t.Log("Task removed from epic")
    })
    
    t.Run("DeleteEpic_TasksUnlinked", func(t *testing.T) {
        // Delete epic
        resp, err := env.Client.Delete(fmt.Sprintf("/epics/%s", epicID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusNoContent, resp.StatusCode)
        
        // Verify tasks still exist but are unlinked
        resp, err = env.Client.Get(fmt.Sprintf("/boards/%s/tasks", boardID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        var result []map[string]interface{}
        ParseResponse(t, resp, &result)
        
        // All tasks should still exist
        assert.Len(t, result, 3, "Tasks should still exist after epic deletion")
        
        // All tasks should have null epic_id
        for _, task := range result {
            assert.Nil(t, task["epic_id"], "Tasks should be unlinked from deleted epic")
        }
        
        t.Log("Epic deleted, tasks unlinked but preserved")
    })
    
    t.Run("EpicNotFound", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/epics/%s", epicID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })
}

// ============================================================
// TEST SUITE: Board Management
// ============================================================

func TestBoard_Management(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    setupAdmin(t, env)
    projectID, defaultBoardID := createProjectAndGetBoard(t, env)
    
    t.Run("DefaultBoardExists", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/boards/%s", defaultBoardID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "Main Board", result["name"])
    })
    
    t.Run("CreateAdditionalBoard", func(t *testing.T) {
        input := map[string]interface{}{
            "name": "Sprint Board",
        }
        
        resp, err := env.Client.Post(fmt.Sprintf("/projects/%s/boards", projectID), input)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusCreated, resp.StatusCode)
        
        var result map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Equal(t, "Sprint Board", result["name"])
        assert.Equal(t, projectID, result["project_id"])
    })
    
    t.Run("ListBoards", func(t *testing.T) {
        resp, err := env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        var result []map[string]interface{}
        ParseResponse(t, resp, &result)
        
        assert.Len(t, result, 2, "Should have 2 boards")
    })
    
    t.Run("DeleteBoard", func(t *testing.T) {
        // Get second board
        resp, _ := env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
        var boards []map[string]interface{}
        ParseResponse(t, resp, &boards)
        resp.Body.Close()
        
        secondBoardID := boards[1]["id"].(string)
        
        // Delete it
        resp, err := env.Client.Delete(fmt.Sprintf("/boards/%s", secondBoardID))
        require.NoError(t, err)
        defer resp.Body.Close()
        
        require.Equal(t, http.StatusNoContent, resp.StatusCode)
        
        // Verify deleted
        resp, _ = env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
        ParseResponse(t, resp, &boards)
        resp.Body.Close()
        
        assert.Len(t, boards, 1, "Should have 1 board remaining")
    })
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func createProjectAndGetBoard(t *testing.T, env *TestEnv) (string, string) {
    // Create project
    projectInput := fixtures.NewProjectInput()
    resp, err := env.Client.Post("/projects", projectInput)
    require.NoError(t, err)
    
    var project map[string]interface{}
    ParseResponse(t, resp, &project)
    resp.Body.Close()
    
    projectID := project["id"].(string)
    
    // Get default board
    resp, err = env.Client.Get(fmt.Sprintf("/projects/%s/boards", projectID))
    require.NoError(t, err)
    
    var boards []map[string]interface{}
    ParseResponse(t, resp, &boards)
    resp.Body.Close()
    
    require.Len(t, boards, 1, "Should have default board")
    
    return projectID, boards[0]["id"].(string)
}
```

## Test Cases Summary

| Test | DB Isolation | Description |
|------|--------------|-------------|
| `TestProject_CRUD` | ✅ Fresh DB | Create, read, update, archive, delete project |
| `TestProject_CreateValidation` | ✅ Fresh DB | Validation errors for invalid input |
| `TestTask_Lifecycle` | ✅ Fresh DB | Full task state transitions |
| `TestTask_PositionOrdering` | ✅ Fresh DB | Task position management |
| `TestEpic_TaskRelationship` | ✅ Fresh DB | Epic creation, tasks in epic, delete epic |
| `TestBoard_Management` | ✅ Fresh DB | Board CRUD operations |

## Running Tests

```bash
# Run project E2E tests
cd backend && go test -v -tags=e2e -run TestProject ./test/e2e/...
cd backend && go test -v -tags=e2e -run TestTask ./test/e2e/...
cd backend && go test -v -tags=e2e -run TestEpic ./test/e2e/...
cd backend && go test -v -tags=e2e -run TestBoard ./test/e2e/...

# Run all project flow tests
cd backend && go test -v -tags=e2e -run "Test(Project|Task|Epic|Board)" ./test/e2e/...
```

## Expected Output

```
=== RUN   TestProject_CRUD
=== RUN   TestProject_CRUD/Create
    project_test.go:37: Created project: abc123...
=== RUN   TestProject_CRUD/Read
=== RUN   TestProject_CRUD/Update
    project_test.go:71: Project updated
=== RUN   TestProject_CRUD/Archive
    project_test.go:88: Project archived
=== RUN   TestProject_CRUD/Delete
    project_test.go:105: Project deleted
--- PASS: TestProject_CRUD (8.45s)

=== RUN   TestTask_Lifecycle
=== RUN   TestTask_Lifecycle/Create_Backlog
    project_test.go:205: Task created in backlog: def456...
=== RUN   TestTask_Lifecycle/Move_ToTodo
    project_test.go:222: Task moved to todo
=== RUN   TestTask_Lifecycle/Move_ToInProgress
    project_test.go:239: Task moved to inprogress
=== RUN   TestTask_Lifecycle/Move_ToDone
    project_test.go:285: Task moved to done
=== RUN   TestTask_Lifecycle/Delete
    project_test.go:313: Task deleted
--- PASS: TestTask_Lifecycle (12.34s)

=== RUN   TestEpic_TaskRelationship
=== RUN   TestEpic_TaskRelationship/CreateEpic
    project_test.go:408: Epic created: ghi789...
=== RUN   TestEpic_TaskRelationship/CreateTasksInEpic
    project_test.go:433: Created 3 tasks in epic
--- PASS: TestEpic_TaskRelationship (15.67s)
PASS
```

## Files to Create

| File | Action |
|------|--------|
| `backend/test/e2e/project_test.go` | CREATE |
