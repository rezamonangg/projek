# Backend Implementation Plan: Tasks Handler

**Priority:** HIGH
**Estimated Time:** 1.5 hours
**Dependencies:** Boards handler

## Current State

### What Exists
- `backend/internal/task/model.go` - Task struct + status types
- `backend/internal/task/repository.go` - Repository interface + implementation
- `backend/internal/task/label_repository.go` - Label repository

### What's Missing
- `backend/internal/task/service.go` - Service layer
- `backend/internal/task/handler.go` - HTTP handlers

## Frontend API Contract

Based on `frontend/src/lib/api/real/tasks.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/boards/{boardId}/tasks` | - | `Task[]` |
| POST | `/boards/{boardId}/tasks` | `{title, description, epic_id?, assignee_id?, reporter_id, status}` | `Task` |
| PUT | `/tasks/{id}` | `{title?, description?, status?, epic_id?, assignee_id?}` | `Task` |
| PATCH | `/tasks/{id}/move` | `{status, position}` | `Task` |
| DELETE | `/tasks/{id}` | - | `void` |

## Implementation Steps

### Step 1: Create Service

**File:** `backend/internal/task/service.go`

```go
package task

type Service struct {
    repo       Repository
    labelRepo  LabelRepository
}

func NewService(repo Repository, labelRepo LabelRepository) *Service {
    return &Service{repo: repo, labelRepo: labelRepo}
}

func (s *Service) Create(ctx context.Context, boardID uuid.UUID, input CreateTaskInput) (*Task, error) {
    task := &Task{
        ID:          uuid.New(),
        BoardID:     boardID,
        Title:       input.Title,
        Description: input.Description,
        Status:      input.Status,
        Position:    0, // Calculate based on existing tasks
        ReporterID:  input.ReporterID,
        CreatedAt:   common.Now(),
        UpdatedAt:   common.Now(),
    }
    
    if err := s.repo.Create(ctx, task); err != nil {
        return nil, err
    }
    return task, nil
}

func (s *Service) GetByBoard(ctx context.Context, boardID uuid.UUID) ([]Task, error)
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateTaskInput) (*Task, error)
func (s *Service) Move(ctx context.Context, id uuid.UUID, status Status, position int) (*Task, error)
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error
```

### Step 2: Create Handler

**File:** `backend/internal/task/handler.go`

```go
package task

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) BoardRoutes() chi.Router {
    // Routes under /boards/{boardId}/tasks
    r := chi.NewRouter()
    r.Get("/", h.ListByBoard)
    r.Post("/", h.Create)
    return r
}

func (h *Handler) Routes() chi.Router {
    // Routes under /tasks/{id}
    r := chi.NewRouter()
    r.Put("/{id}", h.Update)
    r.Patch("/{id}/move", h.Move)
    r.Delete("/{id}", h.Delete)
    return r
}
```

### Step 3: Implement Handlers

```go
func (h *Handler) ListByBoard(w http.ResponseWriter, r *http.Request) {
    boardID := chi.URLParam(r, "boardId")
    tasks, err := h.service.GetByBoard(r.Context(), uuid.MustParse(boardID))
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    common.Success(w, 200, tasks)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    boardID := chi.URLParam(r, "boardId")
    
    var req CreateTaskInput
    if err := common.ParseJSON(r, &req); err != nil { ... }
    if err := common.ValidateStruct(&req); err != nil { ... }
    
    task, err := h.service.Create(r.Context(), uuid.MustParse(boardID), req)
    // ...
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req UpdateTaskInput
    // ...
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req struct {
        Status   Status `json:"status"`
        Position int    `json:"position"`
    }
    // ...
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    // ...
}
```

### Step 4: Register Routes

**File:** `backend/internal/router/router.go`

```go
taskService := task.NewService(taskRepo, taskLabelRepo)
taskHandler := task.NewHandler(taskService)

// Under boards
r.Route("/boards/{boardId}", func(r chi.Router) {
    r.Mount("/tasks", taskHandler.BoardRoutes())
})

// Direct task access
r.Mount("/tasks", taskHandler.Routes())
```

## Model Definition

**File:** `backend/internal/task/model.go`

```go
type Status string

const (
    StatusBacklog    Status = "backlog"
    StatusTodo       Status = "todo"
    StatusInProgress Status = "inprogress"
    StatusCodeReview Status = "codereview"
    StatusInTest     Status = "intest"
    StatusNeedDeploy Status = "needdeploy"
    StatusDone       Status = "done"
)

type Task struct {
    ID          uuid.UUID  `json:"id"`
    BoardID     uuid.UUID  `json:"board_id"`
    EpicID      *uuid.UUID `json:"epic_id,omitempty"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      Status     `json:"status"`
    Position    int        `json:"position"`
    StoryPoints *int       `json:"story_points,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
    ReporterID  uuid.UUID  `json:"reporter_id"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    Labels      []string   `json:"labels,omitempty"`
}

type CreateTaskInput struct {
    Title       string     `json:"title" validate:"required"`
    Description string     `json:"description"`
    EpicID      *uuid.UUID `json:"epic_id"`
    AssigneeID  *uuid.UUID `json:"assignee_id"`
    ReporterID  uuid.UUID  `json:"reporter_id" validate:"required"`
    Status      Status     `json:"status"`
}

type UpdateTaskInput struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      Status     `json:"status"`
    EpicID      *uuid.UUID `json:"epic_id"`
    AssigneeID  *uuid.UUID `json:"assignee_id"`
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/task/handler_test.go`

```go
package task

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/monachy/projek/internal/common"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

var testBoardID = uuid.MustParse("55555555-5555-5555-5555-555555555555")
var testTaskID = uuid.MustParse("66666666-6666-6666-6666-666666666666")
var testReporterID = uuid.MustParse("77777777-7777-7777-7777-777777777777")

func testTask() *Task {
    now := time.Now()
    return &Task{
        ID:          testTaskID,
        BoardID:     testBoardID,
        Title:       "Test Task",
        Description: "Task description",
        Status:      StatusTodo,
        Position:    0,
        ReporterID:  testReporterID,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}

func setupTaskHandler(t *testing.T) (*Handler, *MockRepository) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    return NewHandler(svc), repo
}
```

### Handler Tests

```go
func TestTaskHandler_ListByBoard_Success(t *testing.T) {
    handler, repo := setupTaskHandler(t)
    
    tasks := []Task{*testTask()}
    repo.EXPECT().GetByBoard(gomock.Any(), testBoardID).Return(tasks, nil)
    
    r := chi.NewRouter()
    r.Get("/{boardId}/tasks", handler.ListByBoard)
    
    req := httptest.NewRequest("GET", "/"+testBoardID.String()+"/tasks", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Task
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
}

func TestTaskHandler_Create_Success(t *testing.T) {
    handler, repo := setupTaskHandler(t)
    
    body := `{"title":"New Task","description":"Desc","reporter_id":"77777777-7777-7777-7777-777777777777"}`
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    r := chi.NewRouter()
    r.Post("/{boardId}/tasks", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testBoardID.String()+"/tasks", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestTaskHandler_Create_MissingTitle(t *testing.T) {
    handler, _ := setupTaskHandler(t)
    
    body := `{"description":"Desc"}`
    
    r := chi.NewRouter()
    r.Post("/{boardId}/tasks", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testBoardID.String()+"/tasks", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Update_Success(t *testing.T) {
    handler, repo := setupTaskHandler(t)
    
    existing := testTask()
    repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    body := `{"title":"Updated Task","status":"inprogress"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/"+testTaskID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTaskHandler_Move_Success(t *testing.T) {
    handler, repo := setupTaskHandler(t)
    
    existing := testTask()
    repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    body := `{"status":"done","position":5}`
    r := chi.NewRouter()
    r.Patch("/{id}/move", handler.Move)
    
    req := httptest.NewRequest("PATCH", "/"+testTaskID.String()+"/move", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTaskHandler_Delete_Success(t *testing.T) {
    handler, repo := setupTaskHandler(t)
    
    repo.EXPECT().Delete(gomock.Any(), testTaskID).Return(nil)
    
    r := chi.NewRouter()
    r.Delete("/{id}", handler.Delete)
    
    req := httptest.NewRequest("DELETE", "/"+testTaskID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNoContent, rr.Code)
}
```

### Service Tests

**File:** `backend/internal/task/service_test.go`

```go
func TestService_Create_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    input := CreateTaskInput{
        Title:      "New Task",
        ReporterID: testReporterID,
        Status:     StatusBacklog,
    }
    
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    task, err := svc.Create(context.Background(), testBoardID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "New Task", task.Title)
    assert.Equal(t, testBoardID, task.BoardID)
    assert.Equal(t, StatusBacklog, task.Status)
}

func TestService_Create_CalculatesPosition(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    input := CreateTaskInput{
        Title:      "New Task",
        ReporterID: testReporterID,
        Status:     StatusBacklog,
    }
    
    repo.EXPECT().GetMaxPosition(gomock.Any(), testBoardID).Return(5, nil)
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, t *Task) error {
        assert.Equal(t, 6, t.Position)
        return nil
    })
    
    task, err := svc.Create(context.Background(), testBoardID, input)
    
    require.NoError(t, err)
    assert.NotNil(t, task)
}

func TestService_Move_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    existing := testTask()
    repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    task, err := svc.Move(context.Background(), testTaskID, StatusDone, 10)
    
    require.NoError(t, err)
    assert.Equal(t, StatusDone, task.Status)
    assert.Equal(t, 10, task.Position)
}

func TestService_GetByBoard_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    tasks := []Task{*testTask()}
    repo.EXPECT().GetByBoard(gomock.Any(), testBoardID).Return(tasks, nil)
    
    result, err := svc.GetByBoard(context.Background(), testBoardID)
    
    require.NoError(t, err)
    assert.Len(t, result, 1)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestTaskHandler_ListByBoard_Success` | Returns tasks for board | 200 |
| `TestTaskHandler_Create_Success` | Creates task | 201 |
| `TestTaskHandler_Create_MissingTitle` | Fails without title | 400 |
| `TestTaskHandler_Update_Success` | Updates task | 200 |
| `TestTaskHandler_Move_Success` | Moves task to new status/position | 200 |
| `TestTaskHandler_Delete_Success` | Deletes task | 204 |
| `TestTaskHandler_NotFound` | Task doesn't exist | 404 |

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/task/service.go` | CREATE |
| `backend/internal/task/service_test.go` | CREATE |
| `backend/internal/task/handler.go` | CREATE |
| `backend/internal/task/handler_test.go` | CREATE |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/task/... -v
```
