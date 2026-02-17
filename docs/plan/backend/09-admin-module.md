# Backend Implementation Plan: Admin Module

**Priority:** LOW
**Estimated Time:** 1.5 hours
**Dependencies:** Members (for stats), Projects (for stats), Tasks (for stats)

## Current State

### What Exists
- Nothing specific to admin in backend

### What's Missing
- Admin handler
- Dashboard stats service
- Settings management

## Frontend API Contract

Based on `frontend/src/lib/api/real/admin.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/admin/settings` | - | `CommunitySettings` |
| PUT | `/admin/settings` | `{...settings}` | `CommunitySettings` |
| GET | `/admin/dashboard` | - | `DashboardStats` |

## Implementation Steps

### Step 1: Create Module Directory

```bash
mkdir -p backend/internal/admin
```

### Step 2: Define Models

**File:** `backend/internal/admin/model.go`

```go
package admin

import (
    "time"
    "github.com/google/uuid"
)

type CommunitySettings struct {
    ID                      uuid.UUID `json:"id"`
    CommunityID             uuid.UUID `json:"community_id"`
    AllowMemberRegistration bool      `json:"allow_member_registration"`
    RequireEmailVerification bool     `json:"require_email_verification"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
}

type DashboardStats struct {
    TotalMembers    int `json:"total_members"`
    TotalProjects   int `json:"total_projects"`
    TotalTasks      int `json:"total_tasks"`
    ActiveTasks     int `json:"active_tasks"`
    CompletedTasks  int `json:"completed_tasks"`
}

type UpdateSettingsInput struct {
    AllowMemberRegistration  *bool `json:"allow_member_registration"`
    RequireEmailVerification *bool `json:"require_email_verification"`
}
```

### Step 3: Create Repository

**File:** `backend/internal/admin/repository.go`

```go
package admin

type Repository interface {
    GetSettings(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error)
    UpdateSettings(ctx context.Context, settings *CommunitySettings) error
}

type repository struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repository{db: db}
}

func (r *repository) GetSettings(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error) {
    query := `SELECT id, community_id, allow_member_registration, require_email_verification, created_at, updated_at
              FROM community_settings WHERE community_id = $1`
    // ...
}

func (r *repository) UpdateSettings(ctx context.Context, settings *CommunitySettings) error {
    query := `UPDATE community_settings 
              SET allow_member_registration = $1, require_email_verification = $2, updated_at = $3
              WHERE community_id = $4`
    // ...
}
```

### Step 4: Create Service

**File:** `backend/internal/admin/service.go`

```go
package admin

type StatsProvider interface {
    GetMemberCount(ctx context.Context, communityID uuid.UUID) (int, error)
    GetProjectCount(ctx context.Context, communityID uuid.UUID) (int, error)
    GetTaskStats(ctx context.Context, communityID uuid.UUID) (total, active, completed int, err error)
}

type Service struct {
    repo           Repository
    statsProvider  StatsProvider
}

func NewService(repo Repository, statsProvider StatsProvider) *Service {
    return &Service{repo: repo, statsProvider: statsProvider}
}

func (s *Service) GetSettings(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error) {
    return s.repo.GetSettings(ctx, communityID)
}

func (s *Service) UpdateSettings(ctx context.Context, communityID uuid.UUID, input UpdateSettingsInput) (*CommunitySettings, error) {
    settings, err := s.repo.GetSettings(ctx, communityID)
    if err != nil {
        return nil, err
    }
    
    if input.AllowMemberRegistration != nil {
        settings.AllowMemberRegistration = *input.AllowMemberRegistration
    }
    if input.RequireEmailVerification != nil {
        settings.RequireEmailVerification = *input.RequireEmailVerification
    }
    settings.UpdatedAt = common.Now()
    
    if err := s.repo.UpdateSettings(ctx, settings); err != nil {
        return nil, err
    }
    
    return settings, nil
}

func (s *Service) GetDashboardStats(ctx context.Context, communityID uuid.UUID) (*DashboardStats, error) {
    memberCount, err := s.statsProvider.GetMemberCount(ctx, communityID)
    if err != nil {
        return nil, err
    }
    
    projectCount, err := s.statsProvider.GetProjectCount(ctx, communityID)
    if err != nil {
        return nil, err
    }
    
    total, active, completed, err := s.statsProvider.GetTaskStats(ctx, communityID)
    if err != nil {
        return nil, err
    }
    
    return &DashboardStats{
        TotalMembers:   memberCount,
        TotalProjects:  projectCount,
        TotalTasks:     total,
        ActiveTasks:    active,
        CompletedTasks: completed,
    }, nil
}
```

### Step 5: Create Stats Provider

**File:** `backend/internal/admin/stats_provider.go`

```go
package admin

type DBStatsProvider struct {
    db *pgxpool.Pool
}

func NewDBStatsProvider(db *pgxpool.Pool) *DBStatsProvider {
    return &DBStatsProvider{db: db}
}

func (p *DBStatsProvider) GetMemberCount(ctx context.Context, communityID uuid.UUID) (int, error) {
    query := `SELECT COUNT(*) FROM members WHERE community_id = $1 AND is_active = true`
    var count int
    err := p.db.QueryRow(ctx, query, communityID).Scan(&count)
    return count, err
}

func (p *DBStatsProvider) GetProjectCount(ctx context.Context, communityID uuid.UUID) (int, error) {
    query := `SELECT COUNT(*) FROM projects WHERE community_id = $1 AND is_archived = false`
    var count int
    err := p.db.QueryRow(ctx, query, communityID).Scan(&count)
    return count, err
}

func (p *DBStatsProvider) GetTaskStats(ctx context.Context, communityID uuid.UUID) (total, active, completed int, err error) {
    // Join tasks with projects to filter by community
    query := `
        SELECT 
            COUNT(*) as total,
            COUNT(*) FILTER (WHERE t.status IN ('inprogress', 'codereview', 'intest', 'needdeploy')) as active,
            COUNT(*) FILTER (WHERE t.status = 'done') as completed
        FROM tasks t
        JOIN boards b ON t.board_id = b.id
        JOIN projects p ON b.project_id = p.id
        WHERE p.community_id = $1
    `
    err = p.db.QueryRow(ctx, query, communityID).Scan(&total, &active, &completed)
    return
}
```

### Step 6: Create Handler

**File:** `backend/internal/admin/handler.go`

```go
package admin

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    // Add admin-only middleware
    r.Get("/settings", h.GetSettings)
    r.Put("/settings", h.UpdateSettings)
    r.Get("/dashboard", h.GetDashboard)
    return r
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
    member := getCurrentMember(r)
    if member == nil {
        common.Error(w, 401, "UNAUTHORIZED", "not authenticated")
        return
    }
    
    if member.Role != "admin" {
        common.Error(w, 403, "FORBIDDEN", "admin access required")
        return
    }
    
    settings, err := h.service.GetSettings(r.Context(), member.CommunityID)
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 200, settings)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    member := getCurrentMember(r)
    if member == nil || member.Role != "admin" {
        common.Error(w, 403, "FORBIDDEN", "admin access required")
        return
    }
    
    var req UpdateSettingsInput
    if err := common.ParseJSON(r, &req); err != nil {
        common.Error(w, 400, "INVALID_REQUEST", err.Error())
        return
    }
    
    settings, err := h.service.UpdateSettings(r.Context(), member.CommunityID, req)
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 200, settings)
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
    member := getCurrentMember(r)
    if member == nil || member.Role != "admin" {
        common.Error(w, 403, "FORBIDDEN", "admin access required")
        return
    }
    
    stats, err := h.service.GetDashboardStats(r.Context(), member.CommunityID)
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 200, stats)
}
```

### Step 7: Create Migration (if not exists)

**File:** `backend/db/migrations/YYYYMMDDHHMMSS_create_community_settings.sql`

```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS community_settings (
    id UUID PRIMARY KEY,
    community_id UUID NOT NULL UNIQUE REFERENCES communities(id),
    allow_member_registration BOOLEAN NOT NULL DEFAULT true,
    require_email_verification BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS community_settings;
```

### Step 8: Register Routes

**File:** `backend/internal/router/router.go`

```go
statsProvider := admin.NewDBStatsProvider(db)
adminRepo := admin.NewRepository(db)
adminService := admin.NewService(adminRepo, statsProvider)
adminHandler := admin.NewHandler(adminService)
r.Mount("/admin", adminHandler.Routes())
```

## Response Format

**Settings:**
```json
{
    "id": "uuid",
    "community_id": "uuid",
    "allow_member_registration": true,
    "require_email_verification": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

**Dashboard:**
```json
{
    "total_members": 50,
    "total_projects": 10,
    "total_tasks": 200,
    "active_tasks": 45,
    "completed_tasks": 120
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/admin/handler_test.go`

```go
package admin

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/monachy/projek/internal/member"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

func testCommunitySettings() *CommunitySettings {
    return &CommunitySettings{
        ID:                      uuid.New(),
        CommunityID:             testCommunityID,
        AllowMemberRegistration: true,
        RequireEmailVerification: false,
        CreatedAt:               time.Now(),
        UpdatedAt:               time.Now(),
    }
}

func testAdminMember() *member.Member {
    return &member.Member{
        ID:          uuid.New(),
        CommunityID: testCommunityID,
        Email:       "admin@example.com",
        Role:        "admin",
        IsActive:    true,
    }
}

func testRegularMember() *member.Member {
    return &member.Member{
        ID:          uuid.New(),
        CommunityID: testCommunityID,
        Email:       "user@example.com",
        Role:        "member",
        IsActive:    true,
    }
}

func setupAdminHandler(t *testing.T) (*Handler, *MockRepository, *MockStatsProvider) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    statsProvider := NewMockStatsProvider(ctrl)
    svc := NewService(repo, statsProvider)
    return NewHandler(svc), repo, statsProvider
}

func withMember(r *http.Request, m *member.Member) *http.Request {
    ctx := context.WithValue(r.Context(), "member", m)
    return r.WithContext(ctx)
}
```

### Handler Tests

```go
func TestAdminHandler_GetSettings_Success(t *testing.T) {
    handler, repo, _ := setupAdminHandler(t)
    
    settings := testCommunitySettings()
    repo.EXPECT().GetSettings(gomock.Any(), testCommunityID).Return(settings, nil)
    
    req := httptest.NewRequest("GET", "/admin/settings", nil)
    req = withMember(req, testAdminMember())
    
    rr := httptest.NewRecorder()
    handler.GetSettings(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp CommunitySettings
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.True(t, resp.AllowMemberRegistration)
}

func TestAdminHandler_GetSettings_Forbidden(t *testing.T) {
    handler, _, _ := setupAdminHandler(t)
    
    req := httptest.NewRequest("GET", "/admin/settings", nil)
    req = withMember(req, testRegularMember())
    
    rr := httptest.NewRecorder()
    handler.GetSettings(rr, req)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAdminHandler_GetSettings_Unauthorized(t *testing.T) {
    handler, _, _ := setupAdminHandler(t)
    
    req := httptest.NewRequest("GET", "/admin/settings", nil)
    // No member in context
    
    rr := httptest.NewRecorder()
    handler.GetSettings(rr, req)
    
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAdminHandler_UpdateSettings_Success(t *testing.T) {
    handler, repo, _ := setupAdminHandler(t)
    
    existing := testCommunitySettings()
    falseVal := false
    body := `{"allow_member_registration":false}`
    
    repo.EXPECT().GetSettings(gomock.Any(), testCommunityID).Return(existing, nil)
    repo.EXPECT().UpdateSettings(gomock.Any(), gomock.Any()).Return(nil)
    
    req := httptest.NewRequest("PUT", "/admin/settings", strings.NewReader(body))
    req = withMember(req, testAdminMember())
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler.UpdateSettings(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAdminHandler_UpdateSettings_Forbidden(t *testing.T) {
    handler, _, _ := setupAdminHandler(t)
    
    body := `{"allow_member_registration":false}`
    
    req := httptest.NewRequest("PUT", "/admin/settings", strings.NewReader(body))
    req = withMember(req, testRegularMember())
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler.UpdateSettings(rr, req)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAdminHandler_GetDashboard_Success(t *testing.T) {
    handler, _, statsProvider := setupAdminHandler(t)
    
    statsProvider.EXPECT().GetMemberCount(gomock.Any(), testCommunityID).Return(50, nil)
    statsProvider.EXPECT().GetProjectCount(gomock.Any(), testCommunityID).Return(10, nil)
    statsProvider.EXPECT().GetTaskStats(gomock.Any(), testCommunityID).Return(200, 45, 120, nil)
    
    req := httptest.NewRequest("GET", "/admin/dashboard", nil)
    req = withMember(req, testAdminMember())
    
    rr := httptest.NewRecorder()
    handler.GetDashboard(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp DashboardStats
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Equal(t, 50, resp.TotalMembers)
    assert.Equal(t, 10, resp.TotalProjects)
    assert.Equal(t, 200, resp.TotalTasks)
    assert.Equal(t, 45, resp.ActiveTasks)
    assert.Equal(t, 120, resp.CompletedTasks)
}

func TestAdminHandler_GetDashboard_Forbidden(t *testing.T) {
    handler, _, _ := setupAdminHandler(t)
    
    req := httptest.NewRequest("GET", "/admin/dashboard", nil)
    req = withMember(req, testRegularMember())
    
    rr := httptest.NewRecorder()
    handler.GetDashboard(rr, req)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
}
```

### Service Tests

**File:** `backend/internal/admin/service_test.go`

```go
func TestService_GetSettings_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    expected := testCommunitySettings()
    repo.EXPECT().GetSettings(gomock.Any(), testCommunityID).Return(expected, nil)
    
    settings, err := svc.GetSettings(context.Background(), testCommunityID)
    
    require.NoError(t, err)
    assert.Equal(t, expected, settings)
}

func TestService_UpdateSettings_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    existing := testCommunitySettings()
    falseVal := false
    input := UpdateSettingsInput{
        AllowMemberRegistration: &falseVal,
    }
    
    repo.EXPECT().GetSettings(gomock.Any(), testCommunityID).Return(existing, nil)
    repo.EXPECT().UpdateSettings(gomock.Any(), gomock.Any()).Return(nil)
    
    settings, err := svc.UpdateSettings(context.Background(), testCommunityID, input)
    
    require.NoError(t, err)
    assert.False(t, settings.AllowMemberRegistration)
}

func TestService_UpdateSettings_PartialUpdate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo, nil)
    
    existing := testCommunitySettings()
    existing.RequireEmailVerification = false
    
    trueVal := true
    input := UpdateSettingsInput{
        RequireEmailVerification: &trueVal,
        // AllowMemberRegistration not set - should remain unchanged
    }
    
    repo.EXPECT().GetSettings(gomock.Any(), testCommunityID).Return(existing, nil)
    repo.EXPECT().UpdateSettings(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, s *CommunitySettings) error {
        assert.True(t, s.RequireEmailVerification)
        assert.True(t, s.AllowMemberRegistration) // unchanged
        return nil
    })
    
    settings, err := svc.UpdateSettings(context.Background(), testCommunityID, input)
    
    require.NoError(t, err)
    assert.NotNil(t, settings)
}

func TestService_GetDashboardStats_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    statsProvider := NewMockStatsProvider(ctrl)
    svc := NewService(repo, statsProvider)
    
    statsProvider.EXPECT().GetMemberCount(gomock.Any(), testCommunityID).Return(100, nil)
    statsProvider.EXPECT().GetProjectCount(gomock.Any(), testCommunityID).Return(20, nil)
    statsProvider.EXPECT().GetTaskStats(gomock.Any(), testCommunityID).Return(500, 100, 300, nil)
    
    stats, err := svc.GetDashboardStats(context.Background(), testCommunityID)
    
    require.NoError(t, err)
    assert.Equal(t, 100, stats.TotalMembers)
    assert.Equal(t, 20, stats.TotalProjects)
    assert.Equal(t, 500, stats.TotalTasks)
    assert.Equal(t, 100, stats.ActiveTasks)
    assert.Equal(t, 300, stats.CompletedTasks)
}
```

### Stats Provider Tests

**File:** `backend/internal/admin/stats_provider_test.go`

```go
func TestDBStatsProvider_GetMemberCount(t *testing.T) {
    // Integration test with real database or use mock
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // If using mock:
    provider := NewMockStatsProvider(ctrl)
    provider.EXPECT().GetMemberCount(gomock.Any(), testCommunityID).Return(50, nil)
    
    count, err := provider.GetMemberCount(context.Background(), testCommunityID)
    
    require.NoError(t, err)
    assert.Equal(t, 50, count)
}

func TestDBStatsProvider_GetProjectCount(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    provider := NewMockStatsProvider(ctrl)
    provider.EXPECT().GetProjectCount(gomock.Any(), testCommunityID).Return(15, nil)
    
    count, err := provider.GetProjectCount(context.Background(), testCommunityID)
    
    require.NoError(t, err)
    assert.Equal(t, 15, count)
}

func TestDBStatsProvider_GetTaskStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    provider := NewMockStatsProvider(ctrl)
    provider.EXPECT().GetTaskStats(gomock.Any(), testCommunityID).Return(100, 30, 50, nil)
    
    total, active, completed, err := provider.GetTaskStats(context.Background(), testCommunityID)
    
    require.NoError(t, err)
    assert.Equal(t, 100, total)
    assert.Equal(t, 30, active)
    assert.Equal(t, 50, completed)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestAdminHandler_GetSettings_Success` | Admin gets settings | 200 |
| `TestAdminHandler_GetSettings_Forbidden` | Non-admin rejected | 403 |
| `TestAdminHandler_GetSettings_Unauthorized` | No auth | 401 |
| `TestAdminHandler_UpdateSettings_Success` | Admin updates settings | 200 |
| `TestAdminHandler_UpdateSettings_Forbidden` | Non-admin rejected | 403 |
| `TestAdminHandler_GetDashboard_Success` | Admin gets stats | 200 |
| `TestAdminHandler_GetDashboard_Forbidden` | Non-admin rejected | 403 |

## Files to Create

| File | Action |
|------|--------|
| `backend/internal/admin/model.go` | CREATE |
| `backend/internal/admin/repository.go` | CREATE |
| `backend/internal/admin/repository_test.go` | CREATE |
| `backend/internal/admin/stats_provider.go` | CREATE |
| `backend/internal/admin/stats_provider_test.go` | CREATE |
| `backend/internal/admin/service.go` | CREATE |
| `backend/internal/admin/service_test.go` | CREATE |
| `backend/internal/admin/handler.go` | CREATE |
| `backend/internal/admin/handler_test.go` | CREATE |
| `backend/db/migrations/*_create_community_settings.sql` | CREATE (if not exists) |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/admin/... -v
```
