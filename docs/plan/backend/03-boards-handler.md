# Backend Implementation Plan: Boards Handler

**Priority:** HIGH
**Estimated Time:** 1 hour
**Dependencies:** Projects handler (shared service)

## Current State

### What Exists
- `backend/internal/project/model.go` - Board struct + CreateBoardInput
- `backend/internal/project/service.go` - CreateBoard, GetBoards, DeleteBoard methods
- `backend/internal/project/epic_board_repository.go` - BoardRepository interface + implementation

### What's Missing
- Board-specific handlers
- Route registration

## Frontend API Contract

Based on `frontend/src/lib/api/real/boards.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/boards` | - | `Board[]` |
| POST | `/projects/{projectId}/boards` | `{name}` | `Board` |
| GET | `/boards/{id}` | - | `Board` |

**Note:** Tasks endpoint is under boards: `/boards/{boardId}/tasks` - see Tasks handler plan

## Implementation Steps

### Step 1: Create Board Handler

**File:** `backend/internal/project/board_handler.go`

```go
package project

type BoardHandler struct {
    service *Service
}

func NewBoardHandler(service *Service) *BoardHandler {
    return &BoardHandler{service: service}
}

func (h *BoardHandler) ProjectRoutes() chi.Router {
    // Routes under /projects/{projectId}/boards
    r := chi.NewRouter()
    r.Get("/", h.ListByProject)
    r.Post("/", h.Create)
    return r
}

func (h *BoardHandler) Routes() chi.Router {
    // Routes under /boards/{id}
    r := chi.NewRouter()
    r.Get("/{id}", h.Get)
    r.Delete("/{id}", h.Delete)
    return r
}
```

### Step 2: Implement Handlers

```go
func (h *BoardHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    boards, err := h.service.GetBoards(r.Context(), uuid.MustParse(projectID))
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    common.Success(w, 200, boards)
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
    projectID := chi.URLParam(r, "projectId")
    
    var req CreateBoardInput
    if err := common.ParseJSON(r, &req); err != nil { ... }
    if err := common.ValidateStruct(&req); err != nil { ... }
    
    board, err := h.service.CreateBoard(r.Context(), uuid.MustParse(projectID), req)
    // ...
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    // Need to add GetBoardByID to service if not exists
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    err := h.service.DeleteBoard(r.Context(), uuid.MustParse(id))
    // ...
}
```

### Step 3: Add Missing Service Method

**File:** `backend/internal/project/service.go`

```go
func (s *Service) GetBoardByID(ctx context.Context, id uuid.UUID) (*Board, error) {
    return s.boardRepo.GetByID(ctx, id)
}
```

### Step 4: Register Routes

**File:** `backend/internal/router/router.go`

```go
boardHandler := project.NewBoardHandler(projectService)

// Under projects
r.Route("/projects/{projectId}", func(r chi.Router) {
    r.Mount("/boards", boardHandler.ProjectRoutes())
})

// Direct board access
r.Mount("/boards", boardHandler.Routes())
```

## Response Format

```json
{
    "id": "uuid",
    "project_id": "uuid",
    "name": "Board Name",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/project/board_handler_test.go`

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
    "github.com/monachy/projek/internal/common"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

var testBoardID = uuid.MustParse("55555555-5555-5555-5555-555555555555")

func testBoard() *Board {
    now := time.Now()
    return &Board{
        ID:        testBoardID,
        ProjectID: testProjectID,
        Name:      "Main Board",
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func setupBoardHandler(t *testing.T) (*BoardHandler, *MockBoardRepository) {
    ctrl := gomock.NewController(t)
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(nil, nil, boardRepo, nil, nil)
    return NewBoardHandler(svc), boardRepo
}
```

### Handler Tests

```go
func TestBoardHandler_ListByProject_Success(t *testing.T) {
    handler, boardRepo := setupBoardHandler(t)
    
    boards := []Board{*testBoard()}
    boardRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(boards, nil)
    
    r := chi.NewRouter()
    r.Get("/{projectId}/boards", handler.ListByProject)
    
    req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/boards", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Board
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
    assert.Equal(t, "Main Board", resp[0].Name)
}

func TestBoardHandler_Create_Success(t *testing.T) {
    handler, boardRepo := setupBoardHandler(t)
    
    body := `{"name":"New Board"}`
    boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    r := chi.NewRouter()
    r.Post("/{projectId}/boards", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/boards", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestBoardHandler_Create_MissingName(t *testing.T) {
    handler, _ := setupBoardHandler(t)
    
    body := `{}`
    
    r := chi.NewRouter()
    r.Post("/{projectId}/boards", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/boards", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_Get_Success(t *testing.T) {
    handler, boardRepo := setupBoardHandler(t)
    
    boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(testBoard(), nil)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/"+testBoardID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBoardHandler_Get_NotFound(t *testing.T) {
    handler, boardRepo := setupBoardHandler(t)
    
    boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(nil, common.ErrNotFound)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/"+testBoardID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestBoardHandler_Delete_Success(t *testing.T) {
    handler, boardRepo := setupBoardHandler(t)
    
    boardRepo.EXPECT().Delete(gomock.Any(), testBoardID).Return(nil)
    
    r := chi.NewRouter()
    r.Delete("/{id}", handler.Delete)
    
    req := httptest.NewRequest("DELETE", "/"+testBoardID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNoContent, rr.Code)
}
```

### Service Tests

```go
func TestService_GetBoardByID_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(nil, nil, boardRepo, nil, nil)
    
    expected := testBoard()
    boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(expected, nil)
    
    board, err := svc.GetBoardByID(context.Background(), testBoardID)
    
    require.NoError(t, err)
    assert.Equal(t, expected, board)
}

func TestService_CreateBoard_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(nil, nil, boardRepo, nil, nil)
    
    input := CreateBoardInput{Name: "New Board"}
    boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    board, err := svc.CreateBoard(context.Background(), testProjectID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "New Board", board.Name)
    assert.Equal(t, testProjectID, board.ProjectID)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestBoardHandler_ListByProject_Success` | Returns boards for project | 200 |
| `TestBoardHandler_Create_Success` | Creates board | 201 |
| `TestBoardHandler_Create_MissingName` | Fails without name | 400 |
| `TestBoardHandler_Get_Success` | Returns board by ID | 200 |
| `TestBoardHandler_Get_NotFound` | Board doesn't exist | 404 |
| `TestBoardHandler_Delete_Success` | Deletes board | 204 |

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/project/board_handler.go` | CREATE |
| `backend/internal/project/board_handler_test.go` | CREATE |
| `backend/internal/project/service.go` | MODIFY (add GetBoardByID) |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/project/... -v -run Board
```
