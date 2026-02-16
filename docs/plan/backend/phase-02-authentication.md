# Backend Phase 2: Authentication & Members

Complete authentication system with session-based auth, member management, and invitations.

## Task 2.1: Member Domain Models

**Commits:**
- `feat: add member entity model`
- `feat: add role enum and permissions`
- `feat: add invitation token model`

**Package:** `internal/member/model.go`

**Models:**
```go
type Member struct {
    ID           uuid.UUID
    CommunityID  uuid.UUID
    Email        string
    PasswordHash string
    FirstName    string
    LastName     string
    DateOfBirth  *time.Time
    Role         Role
    IsActive     bool
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Role string
const (
    RoleAdmin  Role = "admin"
    RoleMember Role = "member"
)
```

## Task 2.2: Member Repository

**Commits:**
- `feat: add member repository interface`
- `feat: implement member repository with pgx`
- `feat: add member repository queries`

**Package:** `internal/member/repository.go`

**Methods:**
- Create, GetByID, GetByEmail, Update, Delete
- ListByCommunity, UpdateLastLogin

## Task 2.3: Session Management

**Commits:**
- `feat: add session store interface`
- `feat: implement Redis session store`
- `feat: add session middleware`

**Package:** `internal/auth/session.go`

## Task 2.4: Auth Service

**Commits:**
- `feat: add auth service with login/logout`
- `feat: add password verification`
- `feat: add session creation and validation`

**Package:** `internal/auth/service.go`

## Task 2.5: Auth Handlers

**Commits:**
- `feat: add login handler`
- `feat: add logout handler`
- `feat: add me/current user handler`
- `feat: add auth middleware for protected routes`

**Package:** `internal/auth/handler.go`

## Task 2.6: Member Service

**Commits:**
- `feat: add member service for CRUD operations`
- `feat: add profile update logic`
- `feat: add password change functionality`

**Package:** `internal/member/service.go`

## Task 2.7: Member Handlers

**Commits:**
- `feat: add list members handler`
- `feat: add get member handler`
- `feat: add update member handler`
- `feat: add update profile handler (self-service)`

**Package:** `internal/member/handler.go`

## Task 2.8: Invitation System

**Commits:**
- `feat: add invitation repository`
- `feat: add invitation service`
- `feat: add send invitation handler`
- `feat: add accept invitation handler`

**Package:** `internal/member/invitation.go`

## Task 2.9: Email Service

**Commits:**
- `feat: add email service interface`
- `feat: implement SMTP email provider`
- `feat: add email template renderer`
- `feat: add invitation email template`

**Package:** `internal/email/`

## Task 2.10: Community Registration

**Commits:**
- `feat: add community repository`
- `feat: add community service`
- `feat: add register community handler`

**Package:** `internal/community/`

## Task 2.11: Database Migrations

**Commits:**
- `task: add communities table migration`
- `task: add members table migration`
- `task: add invitations table migration`
- `task: add community_settings table migration`
- `task: add indexes for email lookups`

---

## Total Backend Commits (Phase 2): 27
