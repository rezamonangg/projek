# Backend Implementation Plan: Labels Handler

**Priority:** MEDIUM
**Estimated Time:** 45 minutes
**Dependencies:** Tasks handler

## Current State

### What Exists
- `backend/internal/task/label_repository.go` - LabelRepository interface + implementation
- Labels are referenced in Task model

### What's Missing
- Label model definition (if not in model.go)
- Label service
- Label handler

## Frontend API Contract

Based on `frontend/src/lib/api/real/labels.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/labels` | - | `Label[]` |
| POST | `/projects/{projectId}/labels` | `{name, color}` | `Label` |
| POST | `/tasks/{taskId}/labels` | `{label_id}` | `Task` |

## Implementation Steps

### Step 1: Define Label Model

**File:** `backend/internal/task/model.go` (add)

```go
type Label struct {
    ID        uuid.UUID `json:"id"`
    ProjectID uuid.UUID `json:"project_id"`
    Name      string    `json:"name"`
    Color     string    `json:"color"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateLabelInput struct {
    Name  string `json:"name" validate:"required"`
    Color string `json:"color" validate:"required"`
}

type AssignLabelInput struct {
    LabelID uuid.UUID `json:"label_id" validate:"required"`
}
```

### Step 2: Create Label Service

**File:** `backend/internal/task/label_service.go`

```go
package task

type LabelService struct {
    labelRepo LabelRepository
    taskRepo  Repository
}

func NewLabelService(labelRepo LabelRepository, taskRepo Repository) *LabelService {
    return &LabelService{labelRepo: labelRepo, taskRepo: taskRepo}
}

func (s *LabelService) GetByProject(ctx context.Context, projectID uuid.UUID) ([]Label, error) {
    return s.labelRepo.GetByProject(ctx, projectID)
}

func (s *LabelService) Create(ctx context.Context, projectID uuid.UUID, input CreateLabelInput) (*Label, error) {
    label := &Label{
        ID:        uuid.New(),
        ProjectID: projectID,
        Name:      input.Name,
        Color:     input.Color,
        CreatedAt: common.Now(),
    }
    
    if err := s.labelRepo.Create(ctx, label); err != nil {
        return nil, err
    }
    return label, nil
}

func (s *LabelService) AssignToTask(ctx context.Context, taskID, labelID uuid.UUID) (*Task, error) {
    if err := s.labelRepo.AssignToTask(ctx, taskID, labelID); err != nil {
        return nil, err
    }
    return s.taskRepo.GetByID(ctx, taskID)
}
```

### Step 3: Create Label Handler

**File:** `backend/internal/task/label_handler.go`

```go
package task

type LabelHandler struct {
    service *LabelService
}

func NewLabelHandler(service *LabelService) *LabelHandler {
    return &LabelHandler{service: service}
}

func (h *LabelHandler) ProjectRoutes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.ListByProject)
    r.Post("/", h.Create)
    return r
}

func (h *LabelHandler) TaskRoutes() chi.Router {
    r := chi.NewRouter()
    r.Post("/", h.AssignToTask)
    return r
}

func (h *LabelHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    labels, err := h.service.GetByProject(r.Context(), uuid.MustParse(projectID))
    // ...
}

func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    var req CreateLabelInput
    // ...
}

func (h *LabelHandler) AssignToTask(w http.ResponseWriter, r *http.Request) {
    taskID := chi.URLParam(r, "taskId")
    var req AssignLabelInput
    // ...
}
```

### Step 4: Register Routes

**File:** `backend/internal/router/router.go`

```go
labelService := task.NewLabelService(taskLabelRepo, taskRepo)
labelHandler := task.NewLabelHandler(labelService)

// Under projects
r.Route("/projects/{projectId}", func(r chi.Router) {
    r.Mount("/labels", labelHandler.ProjectRoutes())
})

// Under tasks
r.Route("/tasks/{taskId}", func(r chi.Router) {
    r.Mount("/labels", labelHandler.TaskRoutes())
})
```

## Response Format

**Label:**
```json
{
    "id": "uuid",
    "project_id": "uuid",
    "name": "Bug",
    "color": "#FF0000",
    "created_at": "2024-01-01T00:00:00Z"
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/task/label_handler_test.go`

```go
package task

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

var testLabelID = uuid.MustParse("88888888-8888-8888-8888-888888888888")
var testProjectID = uuid.MustParse("22222222-2222-2222-2222-222222222222")

func testLabel() *Label {
    return &Label{
        ID:        testLabelID,
        ProjectID: testProjectID,
        Name:      "Bug",
        Color:     "#FF0000",
        CreatedAt: time.Now(),
    }
}

func setupLabelHandler(t *testing.T) (*LabelHandler, *MockLabelRepository, *MockRepository) {
    ctrl := gomock.NewController(t)
    labelRepo := NewMockLabelRepository(ctrl)
    taskRepo := NewMockRepository(ctrl)
    svc := NewLabelService(labelRepo, taskRepo)
    return NewLabelHandler(svc), labelRepo, taskRepo
}
```

### Handler Tests

```go
func TestLabelHandler_ListByProject_Success(t *testing.T) {
    handler, labelRepo, _ := setupLabelHandler(t)
    
    labels := []Label{*testLabel()}
    labelRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(labels, nil)
    
    r := chi.NewRouter()
    r.Get("/{projectId}/labels", handler.ListByProject)
    
    req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/labels", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Label
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
    assert.Equal(t, "Bug", resp[0].Name)
}

func TestLabelHandler_Create_Success(t *testing.T) {
    handler, labelRepo, _ := setupLabelHandler(t)
    
    body := `{"name":"Feature","color":"#00FF00"}`
    labelRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    r := chi.NewRouter()
    r.Post("/{projectId}/labels", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/labels", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestLabelHandler_Create_MissingFields(t *testing.T) {
    handler, _, _ := setupLabelHandler(t)
    
    tests := []struct {
        name string
        body string
    }{
        {"missing name", `{"color":"#00FF00"}`},
        {"missing color", `{"name":"Feature"}`},
        {"empty body", `{}`},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            r := chi.NewRouter()
            r.Post("/{projectId}/labels", handler.Create)
            
            req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/labels", strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rr := httptest.NewRecorder()
            r.ServeHTTP(rr, req)
            
            assert.Equal(t, http.StatusBadRequest, rr.Code)
        })
    }
}

func TestLabelHandler_AssignToTask_Success(t *testing.T) {
    handler, labelRepo, taskRepo := setupLabelHandler(t)
    
    labelRepo.EXPECT().AssignToTask(gomock.Any(), testTaskID, testLabelID).Return(nil)
    taskRepo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(testTask(), nil)
    
    body := `{"label_id":"88888888-8888-8888-8888-888888888888"}`
    r := chi.NewRouter()
    r.Post("/{taskId}/labels", handler.AssignToTask)
    
    req := httptest.NewRequest("POST", "/"+testTaskID.String()+"/labels", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}
```

### Service Tests

**File:** `backend/internal/task/label_service_test.go`

```go
func TestLabelService_Create_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    labelRepo := NewMockLabelRepository(ctrl)
    svc := NewLabelService(labelRepo, nil)
    
    input := CreateLabelInput{
        Name:  "Bug",
        Color: "#FF0000",
    }
    
    labelRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    label, err := svc.Create(context.Background(), testProjectID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "Bug", label.Name)
    assert.Equal(t, "#FF0000", label.Color)
    assert.Equal(t, testProjectID, label.ProjectID)
}

func TestLabelService_GetByProject_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    labelRepo := NewMockLabelRepository(ctrl)
    svc := NewLabelService(labelRepo, nil)
    
    labels := []Label{*testLabel()}
    labelRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(labels, nil)
    
    result, err := svc.GetByProject(context.Background(), testProjectID)
    
    require.NoError(t, err)
    assert.Len(t, result, 1)
}

func TestLabelService_AssignToTask_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    labelRepo := NewMockLabelRepository(ctrl)
    taskRepo := NewMockRepository(ctrl)
    svc := NewLabelService(labelRepo, taskRepo)
    
    labelRepo.EXPECT().AssignToTask(gomock.Any(), testTaskID, testLabelID).Return(nil)
    taskRepo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(testTask(), nil)
    
    task, err := svc.AssignToTask(context.Background(), testTaskID, testLabelID)
    
    require.NoError(t, err)
    assert.NotNil(t, task)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestLabelHandler_ListByProject_Success` | Returns labels for project | 200 |
| `TestLabelHandler_Create_Success` | Creates label | 201 |
| `TestLabelHandler_Create_MissingFields` | Fails without name/color | 400 |
| `TestLabelHandler_AssignToTask_Success` | Assigns label to task | 200 |

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/task/model.go` | MODIFY (add Label struct) |
| `backend/internal/task/label_service.go` | CREATE |
| `backend/internal/task/label_service_test.go` | CREATE |
| `backend/internal/task/label_handler.go` | CREATE |
| `backend/internal/task/label_handler_test.go` | CREATE |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/task/... -v -run Label
```
