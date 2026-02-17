# Backend Implementation Plan: Projects Handler

**Priority:** HIGH (blocking user flow)
**Estimated Time:** 1-2 hours
**Dependencies:** None (model, service, repository already exist)

## Current State

### What Exists
- `backend/internal/project/model.go` - Project, Epic, Board structs + input types
- `backend/internal/project/service.go` - Full CRUD service methods
- `backend/internal/project/repository.go` - Repository interface + implementation
- `backend/internal/project/epic_board_repository.go` - Epic/Board repositories

### What's Missing
- `backend/internal/project/handler.go` - HTTP handlers
- Route registration in `backend/internal/router/router.go`

## Frontend API Contract

Based on `frontend/src/lib/api/real/projects.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects` | - | `{items: Project[], total, page, page_size, total_pages}` |
| POST | `/projects` | `{name, description, key}` | `Project` |
| GET | `/projects/{id}` | - | `Project` |
| PUT | `/projects/{id}` | `{name?, description?, is_archived?}` | `Project` |
| DELETE | `/projects/{id}` | - | `void` |

## Implementation Steps

### Step 1: Create Handler File

**File:** `backend/internal/project/handler.go`

```go
package project

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Route("/{id}", func(r chi.Router) {
        r.Get("/", h.Get)
        r.Put("/", h.Update)
        r.Delete("/", h.Delete)
    })
    return r
}
```

### Step 2: Implement Handlers

**List Handler:**
```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    // Get community_id from authenticated member context
    member := getCurrentMember(r)
    if member == nil {
        common.Error(w, 401, "UNAUTHORIZED", "not authenticated")
        return
    }
    
    // Parse pagination query params
    page := queryInt(r, "page", 1)
    pageSize := queryInt(r, "page_size", 20)
    offset := (page - 1) * pageSize
    
    projects, err := h.service.GetByCommunity(r.Context(), member.CommunityID, pageSize, offset)
    // ... return paginated response
}
```

**Create Handler:**
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    member := getCurrentMember(r)
    
    var req CreateProjectInput
    if err := common.ParseJSON(r, &req); err != nil { ... }
    if err := common.ValidateStruct(&req); err != nil { ... }
    
    project, err := h.service.Create(r.Context(), member.CommunityID, req)
    // ... return created project
}
```

**Get Handler:**
```go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    project, err := h.service.GetByID(r.Context(), uuid.MustParse(id))
    // ... return project or 404
}
```

**Update Handler:**
```go
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req UpdateProjectInput
    // ... parse, validate, update
}
```

**Delete Handler:**
```go
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    err := h.service.Delete(r.Context(), uuid.MustParse(id))
    // ... return success or error
}
```

### Step 3: Register Routes

**File:** `backend/internal/router/router.go`

```go
// In NewRouter function, after auth setup:

// Project module
projectRepo := project.NewRepository(db)
projectEpicRepo := project.NewEpicRepository(db)
projectBoardRepo := project.NewBoardRepository(db)
taskRepo := task.NewRepository(db)
taskLabelRepo := task.NewLabelRepository(db)

projectService := project.NewService(
    projectRepo,
    projectEpicRepo,
    projectBoardRepo,
    taskRepo,
    taskLabelRepo,
)
projectHandler := project.NewHandler(projectService)
r.Mount("/projects", projectHandler.Routes())
```

### Step 4: Authentication Middleware

The handler needs access to the authenticated member. Options:

1. **Context value** (current pattern): Use `getCurrentMember(r)` helper
2. **Middleware**: Add auth middleware that injects member into context

Verify `auth.handler.go` pattern and replicate.

## Response Format

Backend returns snake_case, frontend transforms to camelCase:

```json
{
    "id": "uuid",
    "community_id": "uuid",
    "name": "Project Name",
    "description": "Description",
    "key": "KEY",
    "is_archived": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

For list endpoint:
```json
{
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/project/handler_test.go`

```go
package project

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/monachy/projek/internal/common"
    "github.com/monachy/projek/internal/member"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

// Test helpers
var testCommunityID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testProjectID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
var testMemberID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

func testMember() *member.Member {
    return &member.Member{
        ID:          testMemberID,
        CommunityID: testCommunityID,
        Email:       "test@example.com",
        Role:        "member",
        IsActive:    true,
    }
}

func testAdminMember() *member.Member {
    m := testMember()
    m.Role = "admin"
    return m
}

func testProject() *Project {
    now := time.Now()
    return &Project{
        ID:          testProjectID,
        CommunityID: testCommunityID,
        Name:        "Test Project",
        Description: "Description",
        Key:         "TEST",
        IsArchived:  false,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}

func setupHandler(t *testing.T) (*Handler, *MockRepository, *MockEpicRepository, *MockBoardRepository) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    epicRepo := NewMockEpicRepository(ctrl)
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(repo, epicRepo, boardRepo, nil, nil)
    return NewHandler(svc), repo, epicRepo, boardRepo
}

func withMember(r *http.Request, m *member.Member) *http.Request {
    ctx := context.WithValue(r.Context(), "member", m)
    return r.WithContext(ctx)
}
```

### Handler Tests

```go
func TestHandler_Create_Success(t *testing.T) {
    handler, repo, _, boardRepo := setupHandler(t)
    
    body := `{"name":"New Project","description":"Desc","key":"NEW"}`
    req := httptest.NewRequest("POST", "/projects", strings.NewReader(body))
    req = withMember(req, testMember())
    req.Header.Set("Content-Type", "application/json")
    
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    rr := httptest.NewRecorder()
    handler.Create(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
    
    var resp Project
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Equal(t, "New Project", resp.Name)
    assert.Equal(t, testCommunityID, resp.CommunityID)
}

func TestHandler_Create_ValidationError(t *testing.T) {
    handler, _, _, _ := setupHandler(t)
    
    tests := []struct {
        name string
        body string
    }{
        {"empty name", `{"name":"","description":"Desc","key":"NEW"}`},
        {"missing key", `{"name":"Project","description":"Desc"}`},
        {"key too short", `{"name":"Project","key":"A"}`},
        {"key too long", `{"name":"Project","key":"TOOLONGKEY"}`},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/projects", strings.NewReader(tt.body))
            req = withMember(req, testMember())
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            handler.Create(rr, req)
            
            assert.Equal(t, http.StatusBadRequest, rr.Code)
        })
    }
}

func TestHandler_Create_Unauthorized(t *testing.T) {
    handler, _, _, _ := setupHandler(t)
    
    body := `{"name":"New Project","description":"Desc","key":"NEW"}`
    req := httptest.NewRequest("POST", "/projects", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    // No member in context
    
    rr := httptest.NewRecorder()
    handler.Create(rr, req)
    
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandler_List_Success(t *testing.T) {
    handler, repo, _, _ := setupHandler(t)
    
    projects := []Project{*testProject()}
    repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID, 20, 0).Return(projects, nil)
    
    req := httptest.NewRequest("GET", "/projects", nil)
    req = withMember(req, testMember())
    
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp PaginatedResponse[Project]
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp.Items, 1)
    assert.Equal(t, "Test Project", resp.Items[0].Name)
}

func TestHandler_Get_Success(t *testing.T) {
    handler, repo, _, _ := setupHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(testProject(), nil)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/projects/"+testProjectID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
    handler, repo, _, _ := setupHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(nil, common.ErrNotFound)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/projects/"+testProjectID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_Update_Success(t *testing.T) {
    handler, repo, _, _ := setupHandler(t)
    
    existing := testProject()
    repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *Project) error {
        assert.Equal(t, "Updated Name", p.Name)
        return nil
    })
    
    body := `{"name":"Updated Name"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/projects/"+testProjectID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
    handler, repo, _, _ := setupHandler(t)
    
    repo.EXPECT().Delete(gomock.Any(), testProjectID).Return(nil)
    
    r := chi.NewRouter()
    r.Delete("/{id}", handler.Delete)
    
    req := httptest.NewRequest("DELETE", "/projects/"+testProjectID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNoContent, rr.Code)
}
```

### Service Tests (if not exists)

**File:** `backend/internal/project/service_test.go`

```go
func TestService_Create_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(repo, nil, boardRepo, nil, nil)
    
    input := CreateProjectInput{
        Name:        "Test Project",
        Description: "Description",
        Key:         "TEST",
    }
    
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    project, err := svc.Create(context.Background(), testCommunityID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "Test Project", project.Name)
    assert.Equal(t, testCommunityID, project.CommunityID)
    assert.False(t, project.IsArchived)
}

func TestService_Create_CreatesDefaultBoard(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    boardRepo := NewMockBoardRepository(ctrl)
    svc := NewService(repo, nil, boardRepo, nil, nil)
    
    input := testCreateProjectInput()
    
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, b *Board) error {
        assert.Equal(t, "Main Board", b.Name)
        return nil
    })
    
    project, err := svc.Create(context.Background(), testCommunityID, input)
    
    require.NoError(t, err)
    assert.NotNil(t, project)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestHandler_Create_Success` | Valid input creates project | 201 |
| `TestHandler_Create_ValidationError` | Missing/invalid fields | 400 |
| `TestHandler_Create_Unauthorized` | No auth member in context | 401 |
| `TestHandler_Create_RepoError` | Database error | 500 |
| `TestHandler_List_Success` | Returns paginated projects | 200 |
| `TestHandler_List_Unauthorized` | No auth member | 401 |
| `TestHandler_Get_Success` | Returns project by ID | 200 |
| `TestHandler_Get_NotFound` | Project doesn't exist | 404 |
| `TestHandler_Update_Success` | Partial update works | 200 |
| `TestHandler_Update_NotFound` | Project doesn't exist | 404 |
| `TestHandler_Delete_Success` | Delete succeeds | 204 |
| `TestHandler_Delete_NotFound` | Project doesn't exist | 404 |

## Verification Commands

```bash
# Build
cd backend && go build ./...

# Test
cd backend && go test ./internal/project/... -v

# Test with coverage
cd backend && go test ./internal/project/... -coverprofile=coverage.out
cd backend && go tool cover -html=coverage.out

# Lint
cd backend && golangci-lint run ./internal/project/...
```

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/project/handler.go` | CREATE |
| `backend/internal/project/handler_test.go` | CREATE |
| `backend/internal/project/service_test.go` | CREATE (if not exists) |
| `backend/internal/router/router.go` | MODIFY (add route mount) |
