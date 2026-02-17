# Backend Implementation Plan: Epics Handler

**Priority:** HIGH
**Estimated Time:** 1 hour
**Dependencies:** Projects handler (shared service)

## Current State

### What Exists
- `backend/internal/project/model.go` - Epic struct + CreateEpicInput
- `backend/internal/project/service.go` - CreateEpic, GetEpics, UpdateEpic, DeleteEpic methods
- `backend/internal/project/epic_board_repository.go` - EpicRepository interface + implementation

### What's Missing
- Epic-specific handlers
- Route registration for epic endpoints

## Frontend API Contract

Based on `frontend/src/lib/api/real/epics.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/epics` | - | `Epic[]` |
| POST | `/projects/{projectId}/epics` | `{name, description}` | `Epic` |
| PUT | `/epics/{id}` | `{name, description}` | `Epic` |
| DELETE | `/epics/{id}` | - | `void` |

## Implementation Options

### Option A: Separate Epic Handler (Recommended)
- Clean separation of concerns
- Follows pattern of `/epics/{id}` for direct access

### Option B: Nested in Project Handler
- All routes under `/projects/{id}/epics`
- Simpler but less flexible

**Recommendation:** Option A - separate handler for `/epics/{id}` routes

## Implementation Steps

### Step 1: Create Epic Handler

**File:** `backend/internal/project/epic_handler.go`

```go
package project

type EpicHandler struct {
    service *Service
}

func NewEpicHandler(service *Service) *EpicHandler {
    return &EpicHandler{service: service}
}

func (h *EpicHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Route("/{id}", func(r chi.Router) {
        r.Put("/", h.Update)
        r.Delete("/", h.Delete)
    })
    return r
}
```

### Step 2: Implement Handlers

```go
func (h *EpicHandler) List(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    epics, err := h.service.GetEpics(r.Context(), uuid.MustParse(projectID))
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    common.Success(w, 200, epics)
}

func (h *EpicHandler) Create(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    
    var req CreateEpicInput
    if err := common.ParseJSON(r, &req); err != nil { ... }
    if err := common.ValidateStruct(&req); err != nil { ... }
    
    epic, err := h.service.CreateEpic(r.Context(), uuid.MustParse(projectID), req)
    // ...
}

func (h *EpicHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req CreateEpicInput  // Reuse for update
    // ...
}

func (h *EpicHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    err := h.service.DeleteEpic(r.Context(), uuid.MustParse(id))
    // ...
}
```

### Step 3: Register Routes

**File:** `backend/internal/router/router.go`

```go
// After project handler setup:
epicHandler := project.NewEpicHandler(projectService)

// Mount epic routes under projects for list/create
r.Route("/projects/{projectId}", func(r chi.Router) {
    r.Mount("/epics", epicHandler.Routes())
})

// Mount direct epic access for update/delete
r.Mount("/epics", epicHandler.DirectRoutes()) // Or handle in same Routes()
```

**Alternative approach:**
Mount all epic routes on `/epics` and require `project_id` in request body or query for create/list.

## Response Format

```json
{
    "id": "uuid",
    "project_id": "uuid",
    "name": "Epic Name",
    "description": "Description",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-02-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/project/epic_handler_test.go`

```go
package project

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

var testEpicID = uuid.MustParse("44444444-4444-4444-4444-444444444444")

func testEpic() *Epic {
    now := time.Now()
    startDate := now.Add(-24 * time.Hour)
    endDate := now.Add(7 * 24 * time.Hour)
    return &Epic{
        ID:          testEpicID,
        ProjectID:   testProjectID,
        Name:        "Test Epic",
        Description: "Epic description",
        StartDate:   &startDate,
        EndDate:     &endDate,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}

func setupEpicHandler(t *testing.T) (*EpicHandler, *MockEpicRepository) {
    ctrl := gomock.NewController(t)
    epicRepo := NewMockEpicRepository(ctrl)
    svc := NewService(nil, epicRepo, nil, nil, nil)
    return NewEpicHandler(svc), epicRepo
}
```

### Handler Tests

```go
func TestEpicHandler_List_Success(t *testing.T) {
    handler, epicRepo := setupEpicHandler(t)
    
    epics := []Epic{*testEpic()}
    epicRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(epics, nil)
    
    r := chi.NewRouter()
    r.Get("/{projectId}/epics", handler.List)
    
    req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/epics", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Epic
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
    assert.Equal(t, "Test Epic", resp[0].Name)
}

func TestEpicHandler_Create_Success(t *testing.T) {
    handler, epicRepo := setupEpicHandler(t)
    
    body := `{"name":"New Epic","description":"Description"}`
    epicRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    r := chi.NewRouter()
    r.Post("/{projectId}/epics", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/epics", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestEpicHandler_Create_MissingName(t *testing.T) {
    handler, _ := setupEpicHandler(t)
    
    body := `{"description":"Description"}`
    
    r := chi.NewRouter()
    r.Post("/{projectId}/epics", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/epics", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_Update_Success(t *testing.T) {
    handler, epicRepo := setupEpicHandler(t)
    
    existing := testEpic()
    epicRepo.EXPECT().GetByID(gomock.Any(), testEpicID).Return(existing, nil)
    epicRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    body := `{"name":"Updated Epic","description":"Updated"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/"+testEpicID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEpicHandler_Delete_Success(t *testing.T) {
    handler, epicRepo := setupEpicHandler(t)
    
    epicRepo.EXPECT().Delete(gomock.Any(), testEpicID).Return(nil)
    
    r := chi.NewRouter()
    r.Delete("/{id}", handler.Delete)
    
    req := httptest.NewRequest("DELETE", "/"+testEpicID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNoContent, rr.Code)
}
```

### Service Tests

```go
func TestService_CreateEpic_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    epicRepo := NewMockEpicRepository(ctrl)
    svc := NewService(nil, epicRepo, nil, nil, nil)
    
    input := CreateEpicInput{
        Name:        "New Epic",
        Description: "Description",
    }
    
    epicRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    epic, err := svc.CreateEpic(context.Background(), testProjectID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "New Epic", epic.Name)
    assert.Equal(t, testProjectID, epic.ProjectID)
}

func TestService_GetEpics_EmptyList(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    epicRepo := NewMockEpicRepository(ctrl)
    svc := NewService(nil, epicRepo, nil, nil, nil)
    
    epicRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return([]Epic{}, nil)
    
    epics, err := svc.GetEpics(context.Background(), testProjectID)
    
    require.NoError(t, err)
    assert.Empty(t, epics)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestEpicHandler_List_Success` | Returns epics for project | 200 |
| `TestEpicHandler_Create_Success` | Creates epic with valid input | 201 |
| `TestEpicHandler_Create_MissingName` | Fails without name | 400 |
| `TestEpicHandler_Update_Success` | Updates epic | 200 |
| `TestEpicHandler_Update_NotFound` | Epic doesn't exist | 404 |
| `TestEpicHandler_Delete_Success` | Deletes epic | 204 |
| `TestEpicHandler_Delete_NotFound` | Epic doesn't exist | 404 |

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/project/epic_handler.go` | CREATE |
| `backend/internal/project/epic_handler_test.go` | CREATE |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/project/... -v -run Epic
```
