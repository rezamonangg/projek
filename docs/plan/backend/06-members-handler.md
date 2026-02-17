# Backend Implementation Plan: Members Handler

**Priority:** MEDIUM
**Estimated Time:** 1 hour
**Dependencies:** None

## Current State

### What Exists
- `backend/internal/member/model.go` - Member struct
- `backend/internal/member/service.go` - Member service
- `backend/internal/member/repository.go` - Repository interface + implementation
- `backend/internal/member/invitation.go` - Invitation logic

### What's Missing
- `backend/internal/member/handler.go` - HTTP handlers
- Route registration

## Frontend API Contract

Based on `frontend/src/lib/api/real/members.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/members` | - | `{items: Member[], total, page, page_size, total_pages}` |
| POST | `/members/invite` | `{email, role}` | `{id: invitationId}` |
| PUT | `/members/{id}` | `{first_name, last_name}` | `Member` |

## Implementation Steps

### Step 1: Create Handler

**File:** `backend/internal/member/handler.go`

```go
package member

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/invite", h.Invite)
    r.Put("/{id}", h.Update)
    return r
}
```

### Step 2: Implement Handlers

```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    // Get member from context (need auth)
    member := getCurrentMember(r)
    if member == nil {
        common.Error(w, 401, "UNAUTHORIZED", "not authenticated")
        return
    }
    
    page := queryInt(r, "page", 1)
    pageSize := queryInt(r, "page_size", 20)
    
    members, total, err := h.service.GetByCommunity(
        r.Context(), 
        member.CommunityID,
        pageSize,
        (page-1)*pageSize,
    )
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 200, PaginatedResponse{
        Items:      members,
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: (total + pageSize - 1) / pageSize,
    })
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
    member := getCurrentMember(r)
    
    // Check if admin
    if member.Role != "admin" {
        common.Error(w, 403, "FORBIDDEN", "only admins can invite members")
        return
    }
    
    var req InviteRequest
    if err := common.ParseJSON(r, &req); err != nil { ... }
    if err := common.ValidateStruct(&req); err != nil { ... }
    
    invitation, err := h.service.Invite(r.Context(), member.CommunityID, req.Email, req.Role)
    if err != nil {
        // Handle duplicate email
        if strings.Contains(err.Error(), "already exists") {
            common.Error(w, 409, "CONFLICT", err.Error())
            return
        }
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 201, map[string]string{"id": invitation.ID.String()})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    var req UpdateProfileRequest
    if err := common.ParseJSON(r, &req); err != nil { ... }
    
    member, err := h.service.UpdateProfile(r.Context(), uuid.MustParse(id), req.FirstName, req.LastName)
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 200, member)
}
```

### Step 3: Add Missing Service Methods (if needed)

Check `backend/internal/member/service.go` for:
- `GetByCommunity(ctx, communityID, limit, offset) ([]Member, int, error)`
- `UpdateProfile(ctx, id, firstName, lastName) (*Member, error)`

Add if not present.

### Step 4: Register Routes

**File:** `backend/internal/router/router.go`

```go
memberHandler := member.NewHandler(memberService)
r.Mount("/members", memberHandler.Routes())
```

## Request/Response Types

```go
type InviteRequest struct {
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role" validate:"required,oneof=admin member"`
}

type UpdateProfileRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}
```

## Response Format

**Member:**
```json
{
    "id": "uuid",
    "community_id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "member",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

**Paginated list:**
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

**File:** `backend/internal/member/handler_test.go`

```go
package member

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/monachy/projek/internal/common"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

func setupMemberHandler(t *testing.T) (*Handler, *MockRepository) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    return NewHandler(svc), repo
}

func withAdminMember(r *http.Request) *http.Request {
    admin := testMember()
    admin.Role = RoleAdmin
    return withMember(r, admin)
}

func withMember(r *http.Request, m *Member) *http.Request {
    ctx := context.WithValue(r.Context(), "member", m)
    return r.WithContext(ctx)
}
```

### Handler Tests

```go
func TestMemberHandler_List_Success(t *testing.T) {
    handler, repo := setupMemberHandler(t)
    
    members := []Member{*testMember()}
    repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID, 20, 0).Return(members, nil)
    
    req := httptest.NewRequest("GET", "/members", nil)
    req = withMember(req, testMember())
    
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp PaginatedResponse[Member]
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp.Items, 1)
}

func TestMemberHandler_List_Unauthorized(t *testing.T) {
    handler, _ := setupMemberHandler(t)
    
    req := httptest.NewRequest("GET", "/members", nil)
    // No member in context
    
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMemberHandler_Invite_Success(t *testing.T) {
    handler, repo := setupMemberHandler(t)
    
    body := `{"email":"new@example.com","role":"member"}`
    
    repo.EXPECT().GetByEmail(gomock.Any(), "new@example.com").Return(nil, common.ErrNotFound)
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
    req = withAdminMember(req)
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler.Invite(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMemberHandler_Invite_ForbiddenForNonAdmin(t *testing.T) {
    handler, _ := setupMemberHandler(t)
    
    body := `{"email":"new@example.com","role":"member"}`
    
    req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
    req = withMember(req, testMember()) // Regular member, not admin
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler.Invite(rr, req)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestMemberHandler_Invite_DuplicateEmail(t *testing.T) {
    handler, repo := setupMemberHandler(t)
    
    body := `{"email":"existing@example.com","role":"member"}`
    
    repo.EXPECT().GetByEmail(gomock.Any(), "existing@example.com").Return(testMember(), nil)
    
    req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
    req = withAdminMember(req)
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler.Invite(rr, req)
    
    assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestMemberHandler_Update_Success(t *testing.T) {
    handler, repo := setupMemberHandler(t)
    
    existing := testMember()
    repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(existing, nil)
    repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
    
    body := `{"first_name":"Updated","last_name":"Name"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/"+testMemberID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMemberHandler_Update_NotFound(t *testing.T) {
    handler, repo := setupMemberHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(nil, common.ErrNotFound)
    
    body := `{"first_name":"Updated"}`
    r := chi.NewRouter()
    r.Put("/{id}", handler.Update)
    
    req := httptest.NewRequest("PUT", "/"+testMemberID.String(), strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
}
```

### Service Tests (additions to existing service_test.go)

```go
func TestService_GetByCommunity_Paginated(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    svc := NewService(repo)
    
    members := []Member{*testMember()}
    repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID, 10, 0).Return(members, nil)
    
    result, err := svc.GetByCommunity(context.Background(), testCommunityID, 10, 0)
    
    require.NoError(t, err)
    assert.Len(t, result, 1)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestMemberHandler_List_Success` | Returns paginated members | 200 |
| `TestMemberHandler_List_Unauthorized` | No auth | 401 |
| `TestMemberHandler_Invite_Success` | Admin invites member | 201 |
| `TestMemberHandler_Invite_Forbidden` | Non-admin tries invite | 403 |
| `TestMemberHandler_Invite_DuplicateEmail` | Email already exists | 409 |
| `TestMemberHandler_Update_Success` | Updates profile | 200 |
| `TestMemberHandler_Update_NotFound` | Member doesn't exist | 404 |

## Files to Create/Modify

| File | Action |
|------|--------|
| `backend/internal/member/handler.go` | CREATE |
| `backend/internal/member/handler_test.go` | CREATE |
| `backend/internal/member/service.go` | MODIFY (if methods missing) |
| `backend/internal/router/router.go` | MODIFY |

## Verification

```bash
cd backend && go test ./internal/member/... -v
```
