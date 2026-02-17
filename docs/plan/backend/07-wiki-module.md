# Backend Implementation Plan: Wiki Module

**Priority:** LOW
**Estimated Time:** 2 hours
**Dependencies:** None

## Current State

### What Exists
- Nothing in backend

### What's Missing
- Model definitions
- Repository + interface
- Service
- Handler
- Database migration

## Frontend API Contract

Based on `frontend/src/lib/api/real/wiki.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/wiki` | - | `WikiPage[]` |
| POST | `/projects/{projectId}/wiki` | `{title, content, parent_id?}` | `WikiPage` |
| GET | `/wiki/{id}` | - | `WikiPage` |
| PUT | `/wiki/{id}` | `{title?, content?, parent_id?}` | `WikiPage` |
| DELETE | `/wiki/{id}` | - | `void` |

## Implementation Steps

### Step 1: Create Module Directory

```bash
mkdir -p backend/internal/wiki
```

### Step 2: Define Model

**File:** `backend/internal/wiki/model.go`

```go
package wiki

import (
    "time"
    "github.com/google/uuid"
)

type WikiPage struct {
    ID        uuid.UUID  `json:"id"`
    ProjectID uuid.UUID  `json:"project_id"`
    ParentID  *uuid.UUID `json:"parent_id,omitempty"`
    Title     string     `json:"title"`
    Slug      string     `json:"slug"`
    Content   string     `json:"content"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type CreateWikiPageInput struct {
    Title    string     `json:"title" validate:"required"`
    Content  string     `json:"content"`
    ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateWikiPageInput struct {
    Title    string     `json:"title"`
    Content  string     `json:"content"`
    ParentID *uuid.UUID `json:"parent_id"`
}
```

### Step 3: Create Repository

**File:** `backend/internal/wiki/repository.go`

```go
package wiki

type Repository interface {
    Create(ctx context.Context, page *WikiPage) error
    GetByID(ctx context.Context, id uuid.UUID) (*WikiPage, error)
    GetByProject(ctx context.Context, projectID uuid.UUID) ([]WikiPage, error)
    Update(ctx context.Context, page *WikiPage) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repository{db: db}
}

// Implement all methods...
```

### Step 4: Create Service

**File:** `backend/internal/wiki/service.go`

```go
package wiki

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID uuid.UUID, input CreateWikiPageInput) (*WikiPage, error) {
    slug := generateSlug(input.Title)
    
    page := &WikiPage{
        ID:        uuid.New(),
        ProjectID: projectID,
        ParentID:  input.ParentID,
        Title:     input.Title,
        Slug:      slug,
        Content:   input.Content,
        CreatedAt: common.Now(),
        UpdatedAt: common.Now(),
    }
    
    if err := s.repo.Create(ctx, page); err != nil {
        return nil, err
    }
    return page, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*WikiPage, error)
func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]WikiPage, error)
func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateWikiPageInput) (*WikiPage, error)
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error

func generateSlug(title string) string {
    // Convert to lowercase, replace spaces with hyphens, remove special chars
    // ...
}
```

### Step 5: Create Handler

**File:** `backend/internal/wiki/handler.go`

```go
package wiki

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) ProjectRoutes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.ListByProject)
    r.Post("/", h.Create)
    return r
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/{id}", h.Get)
    r.Put("/{id}", h.Update)
    r.Delete("/{id}", h.Delete)
    return r
}

// Implement handlers...
```

### Step 6: Create Migration

**File:** `backend/db/migrations/YYYYMMDDHHMMSS_create_wiki_pages.sql`

```sql
-- migrate:up
CREATE TABLE wiki_pages (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES wiki_pages(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(project_id, slug)
);

CREATE INDEX idx_wiki_pages_project_id ON wiki_pages(project_id);
CREATE INDEX idx_wiki_pages_parent_id ON wiki_pages(parent_id);

-- migrate:down
DROP TABLE IF EXISTS wiki_pages;
```

### Step 7: Register Routes

**File:** `backend/internal/router/router.go`

```go
wikiRepo := wiki.NewRepository(db)
wikiService := wiki.NewService(wikiRepo)
wikiHandler := wiki.NewHandler(wikiService)

r.Route("/projects/{projectId}", func(r chi.Router) {
    r.Mount("/wiki", wikiHandler.ProjectRoutes())
})
r.Mount("/wiki", wikiHandler.Routes())
```

## Response Format

```json
{
    "id": "uuid",
    "project_id": "uuid",
    "parent_id": "uuid or null",
    "title": "Page Title",
    "slug": "page-title",
    "content": "{\"type\":\"doc\",\"content\":[...]}",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/wiki/handler_test.go`

```go
package wiki

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

var testWikiPageID = uuid.MustParse("99999999-9999-9999-9999-999999999999")

func testWikiPage() *WikiPage {
    now := time.Now()
    return &WikiPage{
        ID:        testWikiPageID,
        ProjectID: testProjectID,
        Title:     "Test Page",
        Slug:      "test-page",
        Content:   `{"type":"doc","content":[]}`,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func setupWikiHandler(t *testing.T) (*Handler, *MockRepository) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    return NewHandler(svc), repo
}
```

### Handler Tests

```go
func TestWikiHandler_ListByProject_Success(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    pages := []WikiPage{*testWikiPage()}
    repo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(pages, nil)
    
    r := chi.NewRouter()
    r.Get("/{projectId}/wiki", handler.ListByProject)
    
    req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/wiki", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []WikiPage
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
}

func TestWikiHandler_Create_Success(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    body := `{"title":"New Page","content":"{\"type\":\"doc\"}"}`
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    r := chi.NewRouter()
    r.Post("/{projectId}/wiki", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/wiki", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestWikiHandler_Create_WithParent(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    parentID := uuid.New()
    body := `{"title":"Child Page","content":"{}","parent_id":"` + parentID.String() + `"}`
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *WikiPage) error {
        assert.Equal(t, parentID, *p.ParentID)
        return nil
    })
    
    r := chi.NewRouter()
    r.Post("/{projectId}/wiki", handler.Create)
    
    req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/wiki", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestWikiHandler_Get_Success(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testWikiPageID).Return(testWikiPage(), nil)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/"+testWikiPageID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWikiHandler_Get_NotFound(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testWikiPageID).Return(nil, common.ErrNotFound)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Get)
    
    req := httptest.NewRequest("GET", "/"+testWikiPageID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestWikiHandler_Update_Success(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    existing := testWikiPage()
    repo.EXPECT().GetByID(gomock.Any(), testWikiPageID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    body := `{"title":"Updated Title","content":"{\"updated\":true}"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/"+testWikiPageID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWikiHandler_Delete_Success(t *testing.T) {
    handler, repo := setupWikiHandler(t)
    
    repo.EXPECT().Delete(gomock.Any(), testWikiPageID).Return(nil)
    
    r := chi.NewRouter()
    r.Delete("/{id}", handler.Delete)
    
    req := httptest.NewRequest("DELETE", "/"+testWikiPageID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNoContent, rr.Code)
}
```

### Service Tests

**File:** `backend/internal/wiki/service_test.go`

```go
func TestService_Create_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    
    input := CreateWikiPageInput{
        Title:   "New Page",
        Content: `{"type":"doc"}`,
    }
    
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    page, err := svc.Create(context.Background(), testProjectID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "New Page", page.Title)
    assert.Equal(t, "new-page", page.Slug)
    assert.Equal(t, testProjectID, page.ProjectID)
}

func TestService_Create_GeneratesSlug(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    
    tests := []struct {
        title    string
        expected string
    }{
        {"Simple Title", "simple-title"},
        {"Title With  Multiple   Spaces", "title-with-multiple-spaces"},
        {"UPPERCASE TITLE", "uppercase-title"},
        {"Title!@#$%With*Special", "titlewithspecial"},
    }
    
    for _, tt := range tests {
        t.Run(tt.title, func(t *testing.T) {
            input := CreateWikiPageInput{Title: tt.title}
            
            repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *WikiPage) error {
                assert.Equal(t, tt.expected, p.Slug)
                return nil
            })
            
            page, err := svc.Create(context.Background(), testProjectID, input)
            
            require.NoError(t, err)
            assert.NotNil(t, page)
        })
    }
}

func TestService_Update_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    
    existing := testWikiPage()
    repo.EXPECT().GetByID(gomock.Any(), testWikiPageID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    input := UpdateWikiPageInput{
        Title:   "Updated Title",
        Content: `{"updated":true}`,
    }
    
    page, err := svc.Update(context.Background(), testWikiPageID, input)
    
    require.NoError(t, err)
    assert.Equal(t, "Updated Title", page.Title)
}

func TestService_Delete_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    
    repo.EXPECT().Delete(gomock.Any(), testWikiPageID).Return(nil)
    
    err := svc.Delete(context.Background(), testWikiPageID)
    
    require.NoError(t, err)
}
```

### Repository Tests

**File:** `backend/internal/wiki/repository_test.go`

```go
func TestRepository_Create_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    page := testWikiPage()
    
    repo.EXPECT().Create(gomock.Any(), page).Return(nil)
    
    err := repo.Create(context.Background(), page)
    
    require.NoError(t, err)
}

func TestRepository_GetByProject_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    pages := []WikiPage{*testWikiPage()}
    
    repo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(pages, nil)
    
    result, err := repo.GetByProject(context.Background(), testProjectID)
    
    require.NoError(t, err)
    assert.Len(t, result, 1)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestWikiHandler_ListByProject_Success` | Returns wiki pages | 200 |
| `TestWikiHandler_Create_Success` | Creates page | 201 |
| `TestWikiHandler_Create_WithParent` | Creates child page | 201 |
| `TestWikiHandler_Get_Success` | Returns page by ID | 200 |
| `TestWikiHandler_Get_NotFound` | Page doesn't exist | 404 |
| `TestWikiHandler_Update_Success` | Updates page | 200 |
| `TestWikiHandler_Delete_Success` | Deletes page | 204 |

## Files to Create

| File | Action |
|------|--------|
| `backend/internal/wiki/model.go` | CREATE |
| `backend/internal/wiki/repository.go` | CREATE |
| `backend/internal/wiki/repository_test.go` | CREATE |
| `backend/internal/wiki/service.go` | CREATE |
| `backend/internal/wiki/service_test.go` | CREATE |
| `backend/internal/wiki/handler.go` | CREATE |
| `backend/internal/wiki/handler_test.go` | CREATE |
| `backend/db/migrations/*_create_wiki_pages.sql` | CREATE |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/wiki/... -v
```
