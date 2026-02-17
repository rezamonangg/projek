# Phase 2: Authentication & Members

Complete authentication system with session-based auth, member management, and invitations.

## Task 2.1: Member Domain Models

**Commits:**
- `feat: add member entity models`
- `feat: add role enum and permissions`
- `feat: add password reset token model`

**Models:**
```go
type Member struct {
    ID          uuid.UUID
    CommunityID uuid.UUID
    Email       string
    PasswordHash string
    FirstName   string
    LastName    string
    DateOfBirth *time.Time
    Role        Role // admin, member
    IsActive    bool
    LastLoginAt *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Role string
const (
    RoleAdmin  Role = "admin"
    RoleMember Role = "member"
)

type Invitation struct {
    ID          uuid.UUID
    CommunityID uuid.UUID
    Email       string
    Token       string
    ExpiresAt   time.Time
    CreatedBy   uuid.UUID
    CreatedAt   time.Time
}
```

## Task 2.2: Member Repository

**Commits:**
- `feat: add member repository interface`
- `feat: implement member repository with pgx`
- `feat: add member repository queries`

**Methods:**
- Create, GetByID, GetByEmail, Update, Delete
- ListByCommunity, UpdateLastLogin
- GetByInvitationToken

## Task 2.3: Session Management

**Commits:**
- `feat: add session store interface`
- `feat: implement Redis session store`
- `feat: add session middleware`

**Session data:**
```go
type Session struct {
    ID         string
    MemberID   uuid.UUID
    CommunityID uuid.UUID
    Email      string
    Role       Role
    CreatedAt  time.Time
    ExpiresAt  time.Time
}
```

## Task 2.4: Auth Service

**Commits:**
- `feat: add auth service with login/logout`
- `feat: add password verification`
- `feat: add session creation and validation`

## Task 2.5: Auth Handlers

**Commits:**
- `feat: add login handler`
- `feat: add logout handler`
- `feat: add me/current user handler`
- `feat: add auth middleware for protected routes`

## Task 2.6: Member Service

**Commits:**
- `feat: add member service for CRUD operations`
- `feat: add profile update logic`
- `feat: add password change functionality`

## Task 2.7: Member Handlers

**Commits:**
- `feat: add list members handler`
- `feat: add get member handler`
- `feat: add update member handler`
- `feat: add update profile handler (self-service)`

## Task 2.8: Invitation System

**Commits:**
- `feat: add invitation repository`
- `feat: add invitation service`
- `feat: add send invitation handler`
- `feat: add accept invitation handler`

**Flow:**
1. Admin sends invitation (email + token)
2. User clicks link with token
3. User sets password
4. Account activated

## Task 2.9: Email Service

**Commits:**
- `feat: add email service interface`
- `feat: implement SMTP email provider`
- `feat: add email template renderer`
- `feat: add invitation email template`

**Templates:**
- invitation.html
- welcome.html
- password-reset.html

## Task 2.10: Community Registration

**Commits:**
- `feat: add community repository`
- `feat: add community service`
- `feat: add register community handler`

**Flow:**
1. Create community
2. Create admin member
3. Create community settings
4. Send welcome email

## Task 2.11: Database Migrations for Auth

**Commits:**
- `task: add members table migration`
- `task: add invitations table migration`
- `task: add indexes for email lookups`

---

## Total Commits: 27
**Estimated Time:** 3-4 days
