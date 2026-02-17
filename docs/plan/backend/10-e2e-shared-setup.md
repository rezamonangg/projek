# E2E Test Plan: Shared Setup & Infrastructure

**Priority:** HIGH
**Estimated Time:** 1 hour
**Dependencies:** All handlers implemented

## Overview

Shared infrastructure for all E2E tests using testcontainers. Each test suite gets its own isolated PostgreSQL and Redis containers.

## Why Testcontainers over SQLite

| Aspect | Testcontainers | SQLite |
|--------|---------------|--------|
| Real DB behavior | ✅ Uses actual PostgreSQL | ❌ Different SQL dialect |
| Constraints | ✅ Real FK constraints | ❌ Looser constraints |
| JSON operations | ✅ Native JSONB | ❌ Limited JSON support |
| Redis | ✅ Real Redis for sessions | ❌ Not applicable |
| Speed | ❌ 30-60s startup | ✅ Instant |
| Docker required | ✅ Yes | ❌ No |

**Recommendation:** Use testcontainers for E2E (real behavior), SQLite only for unit tests if needed.

## Directory Structure

```
backend/
├── test/
│   ├── e2e/
│   │   ├── setup.go              # Shared testcontainers setup
│   │   ├── client.go             # HTTP test client with cookies
│   │   ├── migrations.go         # Migration runner
│   │   ├── auth_test.go          # Auth flow E2E
│   │   ├── admin_test.go         # Admin flow E2E
│   │   └── project_test.go       # Project flow E2E
│   └── fixtures/
│       └── testdata.go           # Test data builders
└── go.mod
```

## Dependencies

Add to `go.mod`:

```go
require (
    github.com/testcontainers/testcontainers-go v0.29.1
    github.com/testcontainers/testcontainers-go/modules/postgres v0.29.1
    github.com/testcontainers/testcontainers-go/modules/redis v0.29.1
)
```

## Shared Setup File

**File:** `backend/test/e2e/setup.go`

```go
package e2e

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/monachy/projek/internal/common"
    "github.com/monachy/projek/internal/router"
    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
    "github.com/testcontainers/testcontainers-go/wait"
)

// TestEnv holds all resources for a single test suite
type TestEnv struct {
    DB         *pgxpool.Pool
    Redis      *redis.Client
    Server     *httptest.Server
    Client     *TestClient
    CommunityID string
    
    // Container references for cleanup
    pgContainer    *postgres.PostgresContainer
    redisContainer *rediscontainer.RedisContainer
    
    // Cleanup function
    Cleanup func()
}

// SetupTestEnv creates an isolated test environment with fresh containers
// Each call creates a NEW database and redis instance for complete isolation
func SetupTestEnv(t *testing.T) *TestEnv {
    return SetupTestEnvWithSuffix(t, generateRandomSuffix())
}

// SetupTestEnvWithSuffix creates test environment with a specific suffix
// Useful for debugging specific test runs
func SetupTestEnvWithSuffix(t *testing.T, suffix string) *TestEnv {
    ctx := context.Background()
    
    // Start PostgreSQL container with unique database
    dbName := fmt.Sprintf("projek_test_%s", suffix)
    pgContainer, err := postgres.Run(ctx,
        "postgres:15-alpine",
        postgres.WithDatabase(dbName),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(60*time.Second),
        ),
    )
    require.NoError(t, err, "Failed to start PostgreSQL container")
    
    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err, "Failed to get PostgreSQL connection string")
    
    // Connect to PostgreSQL
    db, err := pgxpool.New(ctx, connStr)
    require.NoError(t, err, "Failed to connect to PostgreSQL")
    
    // Verify connection
    err = db.Ping(ctx)
    require.NoError(t, err, "Failed to ping PostgreSQL")
    
    // Run migrations
    err = RunMigrations(ctx, db)
    require.NoError(t, err, "Failed to run migrations")
    
    t.Logf("PostgreSQL ready: %s", dbName)
    
    // Start Redis container
    redisContainer, err := rediscontainer.Run(ctx, "redis:7-alpine")
    require.NoError(t, err, "Failed to start Redis container")
    
    redisHost, err := redisContainer.Host(ctx)
    require.NoError(t, err, "Failed to get Redis host")
    
    redisPort, err := redisContainer.MappedPort(ctx, "6379")
    require.NoError(t, err, "Failed to get Redis port")
    
    redisClient := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
    })
    
    // Verify Redis connection
    err = redisClient.Ping(ctx).Err()
    require.NoError(t, err, "Failed to ping Redis")
    
    t.Logf("Redis ready: %s:%s", redisHost, redisPort.Port())
    
    // Create test server
    cfg := &common.Config{
        SessionSecret: fmt.Sprintf("test-secret-%s", suffix),
        CookieDomain:  "localhost",
    }
    
    logger := zerolog.Nop()
    r := router.NewRouter(cfg, logger, db, redisClient)
    
    server := httptest.NewServer(r)
    t.Logf("Test server ready: %s", server.URL)
    
    // Create cleanup function
    cleanup := func() {
        t.Log("Cleaning up test environment...")
        server.Close()
        db.Close()
        redisClient.Close()
        
        if err := pgContainer.Terminate(ctx); err != nil {
            t.Logf("Warning: failed to terminate PostgreSQL: %v", err)
        }
        if err := redisContainer.Terminate(ctx); err != nil {
            t.Logf("Warning: failed to terminate Redis: %v", err)
        }
        t.Log("Cleanup complete")
    }
    
    return &TestEnv{
        DB:            db,
        Redis:         redisClient,
        Server:        server,
        Client:        NewTestClient(server.URL),
        pgContainer:   pgContainer,
        redisContainer: redisContainer,
        Cleanup:       cleanup,
    }
}

// generateRandomSuffix creates a random suffix for unique database names
func generateRandomSuffix() string {
    bytes := make([]byte, 4)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// SetCommunityID stores the community ID for reference
func (e *TestEnv) SetCommunityID(id string) {
    e.CommunityID = id
}

// QueryDB is a helper for direct database queries in tests
func (e *TestEnv) QueryDB(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
    return e.DB.Query(ctx, query, args...)
}

// QueryRowDB is a helper for single row queries
func (e *TestEnv) QueryRowDB(ctx context.Context, query string, args ...interface{}) pgx.Row {
    return e.DB.QueryRow(ctx, query, args...)
}

// ExecDB is a helper for executing statements
func (e *TestEnv) ExecDB(ctx context.Context, query string, args ...interface{}) error {
    _, err := e.DB.Exec(ctx, query, args...)
    return err
}
```

## Migrations File

**File:** `backend/test/e2e/migrations.go`

```go
package e2e

import (
    "context"
    "embed"
    "fmt"
    "io/fs"
    "sort"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all migrations in order
func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
    // Read all migration files
    entries, err := migrationsFS.ReadDir("migrations")
    if err != nil {
        // If no embedded migrations, use inline
        return runInlineMigrations(ctx, db)
    }
    
    // Sort by filename (should be prefixed with timestamp)
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Name() < entries[j].Name()
    })
    
    for _, entry := range entries {
        if strings.HasSuffix(entry.Name(), ".sql") {
            content, err := migrationsFS.ReadFile("migrations/" + entry.Name())
            if err != nil {
                return fmt.Errorf("read migration %s: %w", entry.Name(), err)
            }
            
            _, err = db.Exec(ctx, string(content))
            if err != nil {
                return fmt.Errorf("execute migration %s: %w", entry.Name(), err)
            }
        }
    }
    
    return nil
}

// runInlineMigrations executes migrations defined inline
// Use this if migrations are not embedded
func runInlineMigrations(ctx context.Context, db *pgxpool.Pool) error {
    migrations := []string{
        migrationCreateCommunities,
        migrationCreateMembers,
        migrationCreateProjects,
        migrationCreateBoards,
        migrationCreateTasks,
        migrationCreateEpics,
        migrationCreateLabels,
        migrationCreateTaskLabels,
        migrationCreateCommunitySettings,
    }
    
    for i, m := range migrations {
        _, err := db.Exec(ctx, m)
        if err != nil {
            return fmt.Errorf("migration %d failed: %w", i, err)
        }
    }
    
    return nil
}

// Inline migrations (simplified versions for testing)

const migrationCreateCommunities = `
CREATE TABLE IF NOT EXISTS communities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

const migrationCreateMembers = `
CREATE TABLE IF NOT EXISTS members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(community_id, email)
);

CREATE INDEX IF NOT EXISTS idx_members_community_id ON members(community_id);
CREATE INDEX IF NOT EXISTS idx_members_email ON members(email);
`

const migrationCreateProjects = `
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    key VARCHAR(10) NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(community_id, key)
);

CREATE INDEX IF NOT EXISTS idx_projects_community_id ON projects(community_id);
`

const migrationCreateBoards = `
CREATE TABLE IF NOT EXISTS boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_boards_project_id ON boards(project_id);
`

const migrationCreateTasks = `
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    epic_id UUID REFERENCES epics(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'backlog',
    position INTEGER NOT NULL DEFAULT 0,
    story_points INTEGER,
    due_date TIMESTAMP WITH TIME ZONE,
    assignee_id UUID REFERENCES members(id) ON DELETE SET NULL,
    reporter_id UUID NOT NULL REFERENCES members(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_board_id ON tasks(board_id);
CREATE INDEX IF NOT EXISTS idx_tasks_epic_id ON tasks(epic_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
`

const migrationCreateEpics = `
CREATE TABLE IF NOT EXISTS epics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_epics_project_id ON epics(project_id);
`

const migrationCreateLabels = `
CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, name)
);

CREATE INDEX IF NOT EXISTS idx_labels_project_id ON labels(project_id);
`

const migrationCreateTaskLabels = `
CREATE TABLE IF NOT EXISTS task_labels (
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);
`

const migrationCreateCommunitySettings = `
CREATE TABLE IF NOT EXISTS community_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL UNIQUE REFERENCES communities(id) ON DELETE CASCADE,
    allow_member_registration BOOLEAN NOT NULL DEFAULT true,
    require_email_verification BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`
```

## Test Client

**File:** `backend/test/e2e/client.go`

```go
package e2e

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "testing"

    "github.com/stretchr/testify/require"
)

// TestClient wraps an HTTP client with cookie support for session handling
type TestClient struct {
    baseURL    string
    httpClient *http.Client
    jar        *cookiejar.Jar
}

// NewTestClient creates a new test client
func NewTestClient(baseURL string) *TestClient {
    jar, _ := cookiejar.New(nil)
    return &TestClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Jar: jar,
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                return http.ErrUseLastResponse // Don't follow redirects
            },
        },
        jar: jar,
    }
}

// Post makes a POST request with JSON body
func (c *TestClient) Post(endpoint string, body interface{}) (*http.Response, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    
    return c.httpClient.Do(req)
}

// Get makes a GET request
func (c *TestClient) Get(endpoint string) (*http.Response, error) {
    req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
    if err != nil {
        return nil, err
    }
    return c.httpClient.Do(req)
}

// Put makes a PUT request with JSON body
func (c *TestClient) Put(endpoint string, body interface{}) (*http.Response, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("PUT", c.baseURL+endpoint, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    
    return c.httpClient.Do(req)
}

// Patch makes a PATCH request with JSON body
func (c *TestClient) Patch(endpoint string, body interface{}) (*http.Response, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("PATCH", c.baseURL+endpoint, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    
    return c.httpClient.Do(req)
}

// Delete makes a DELETE request
func (c *TestClient) Delete(endpoint string) (*http.Response, error) {
    req, err := http.NewRequest("DELETE", c.baseURL+endpoint, nil)
    if err != nil {
        return nil, err
    }
    return c.httpClient.Do(req)
}

// GetCookies returns all cookies for the base URL
func (c *TestClient) GetCookies() []*http.Cookie {
    u, _ := url.Parse(c.baseURL)
    return c.jar.Cookies(u)
}

// GetSessionCookie returns the session_id cookie if present
func (c *TestClient) GetSessionCookie() *http.Cookie {
    for _, cookie := range c.GetCookies() {
        if cookie.Name == "session_id" {
            return cookie
        }
    }
    return nil
}

// ParseResponse parses JSON response into target
func ParseResponse(t *testing.T, resp *http.Response, target interface{}) {
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err, "Failed to read response body")
    
    err = json.Unmarshal(body, target)
    require.NoError(t, err, "Failed to parse JSON response: %s", string(body))
}

// ParseResponseBody returns raw response body as string
func ParseResponseBody(t *testing.T, resp *http.Response) string {
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err, "Failed to read response body")
    
    return string(body)
}

// AssertStatusCode checks response status code
func AssertStatusCode(t *testing.T, expected, actual int, respBody string) {
    if expected != actual {
        t.Errorf("Expected status %d, got %d. Body: %s", expected, actual, respBody)
    }
}

// PrintResponse logs response details for debugging
func PrintResponse(t *testing.T, resp *http.Response, name string) {
    body := ParseResponseBody(t, resp)
    t.Logf("%s: Status=%d, Body=%s", name, resp.StatusCode, body)
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}

// ParseError parses an error response
func ParseError(t *testing.T, resp *http.Response) ErrorResponse {
    var errResp ErrorResponse
    ParseResponse(t, resp, &errResp)
    return errResp
}

// URLBuilder helps construct URLs with path parameters
type URLBuilder struct {
    base string
}

func NewURLBuilder(baseURL string) *URLBuilder {
    return &URLBuilder{base: baseURL}
}

func (b *URLBuilder) Projects() string {
    return "/projects"
}

func (b *URLBuilder) Project(id string) string {
    return fmt.Sprintf("/projects/%s", id)
}

func (b *URLBuilder) ProjectBoards(projectID string) string {
    return fmt.Sprintf("/projects/%s/boards", projectID)
}

func (b *URLBuilder) ProjectEpics(projectID string) string {
    return fmt.Sprintf("/projects/%s/epics", projectID)
}

func (b *URLBuilder) ProjectLabels(projectID string) string {
    return fmt.Sprintf("/projects/%s/labels", projectID)
}

func (b *URLBuilder) Board(boardID string) string {
    return fmt.Sprintf("/boards/%s", boardID)
}

func (b *URLBuilder) BoardTasks(boardID string) string {
    return fmt.Sprintf("/boards/%s/tasks", boardID)
}

func (b *URLBuilder) Task(taskID string) string {
    return fmt.Sprintf("/tasks/%s", taskID)
}

func (b *URLBuilder) TaskMove(taskID string) string {
    return fmt.Sprintf("/tasks/%s/move", taskID)
}

func (b *URLBuilder) TaskLabels(taskID string) string {
    return fmt.Sprintf("/tasks/%s/labels", taskID)
}

func (b *URLBuilder) Epic(epicID string) string {
    return fmt.Sprintf("/epics/%s", epicID)
}

func (b *URLBuilder) Members() string {
    return "/members"
}

func (b *URLBuilder) MemberInvite() string {
    return "/members/invite"
}

func (b *URLBuilder) Member(id string) string {
    return fmt.Sprintf("/members/%s", id)
}

func (b *URLBuilder) AdminDashboard() string {
    return "/admin/dashboard"
}

func (b *URLBuilder) AdminSettings() string {
    return "/admin/settings"
}

func (b *URLBuilder) AuthRegister() string {
    return "/auth/register"
}

func (b *URLBuilder) AuthLogin() string {
    return "/auth/login"
}

func (b *URLBuilder) AuthLogout() string {
    return "/auth/logout"
}

func (b *URLBuilder) AuthMe() string {
    return "/auth/me"
}
```

## Test Data Fixtures

**File:** `backend/test/fixtures/testdata.go`

```go
package fixtures

// RegisterCommunityInput represents community registration data
type RegisterCommunityInput struct {
    Name          string `json:"name"`
    Slug          string `json:"slug"`
    AdminEmail    string `json:"admin_email"`
    AdminPassword string `json:"admin_password"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
}

// NewCommunityInput creates default community registration input
func NewCommunityInput() RegisterCommunityInput {
    return RegisterCommunityInput{
        Name:          "Test Community",
        Slug:          "test-community",
        AdminEmail:    "admin@test.com",
        AdminPassword: "password123",
        FirstName:     "Admin",
        LastName:      "User",
    }
}

// NewCommunityInputWithEmail creates community input with specific email
func NewCommunityInputWithEmail(email string) RegisterCommunityInput {
    input := NewCommunityInput()
    input.AdminEmail = email
    input.Slug = email[:strings.Index(email, "@")]
    return input
}

// LoginInput represents login credentials
type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// NewLoginInput creates login input
func NewLoginInput(email, password string) LoginInput {
    return LoginInput{
        Email:    email,
        Password: password,
    }
}

// CreateProjectInput represents project creation data
type CreateProjectInput struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Key         string `json:"key"`
}

// NewProjectInput creates default project input
func NewProjectInput() CreateProjectInput {
    return CreateProjectInput{
        Name:        "Test Project",
        Description: "A test project for E2E testing",
        Key:         "TEST",
    }
}

// NewProjectInputWithKey creates project input with specific key
func NewProjectInputWithKey(name, key string) CreateProjectInput {
    return CreateProjectInput{
        Name:        name,
        Description: fmt.Sprintf("Project: %s", name),
        Key:         key,
    }
}

// CreateTaskInput represents task creation data
type CreateTaskInput struct {
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Status      string  `json:"status"`
    ReporterID  string  `json:"reporter_id"`
    EpicID      *string `json:"epic_id,omitempty"`
    AssigneeID  *string `json:"assignee_id,omitempty"`
}

// NewTaskInput creates default task input
func NewTaskInput(reporterID string) CreateTaskInput {
    return CreateTaskInput{
        Title:       "Test Task",
        Description: "A test task",
        Status:      "backlog",
        ReporterID:  reporterID,
    }
}

// NewTaskInputWithStatus creates task input with specific status
func NewTaskInputWithStatus(reporterID, title, status string) CreateTaskInput {
    return CreateTaskInput{
        Title:       title,
        Description: fmt.Sprintf("Task: %s", title),
        Status:      status,
        ReporterID:  reporterID,
    }
}

// InviteMemberInput represents member invitation data
type InviteMemberInput struct {
    Email string `json:"email"`
    Role  string `json:"role"`
}

// NewInviteInput creates default invite input
func NewInviteInput(email string) InviteMemberInput {
    return InviteMemberInput{
        Email: email,
        Role:  "member",
    }
}

// NewAdminInviteInput creates admin invite input
func NewAdminInviteInput(email string) InviteMemberInput {
    return InviteMemberInput{
        Email: email,
        Role:  "admin",
    }
}

// CreateLabelInput represents label creation data
type CreateLabelInput struct {
    Name  string `json:"name"`
    Color string `json:"color"`
}

// NewLabelInput creates default label input
func NewLabelInput() CreateLabelInput {
    return CreateLabelInput{
        Name:  "Bug",
        Color: "#FF0000",
    }
}

// NewLabelInputWithValues creates label input with specific values
func NewLabelInputWithValues(name, color string) CreateLabelInput {
    return CreateLabelInput{
        Name:  name,
        Color: color,
    }
}

// CreateEpicInput represents epic creation data
type CreateEpicInput struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    StartDate   *string `json:"start_date,omitempty"`
    EndDate     *string `json:"end_date,omitempty"`
}

// NewEpicInput creates default epic input
func NewEpicInput() CreateEpicInput {
    return CreateEpicInput{
        Name:        "Test Epic",
        Description: "A test epic",
    }
}

// NewEpicInputWithName creates epic input with specific name
func NewEpicInputWithName(name string) CreateEpicInput {
    return CreateEpicInput{
        Name:        name,
        Description: fmt.Sprintf("Epic: %s", name),
    }
}

// UpdateProfileInput represents profile update data
type UpdateProfileInput struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

// NewUpdateProfileInput creates profile update input
func NewUpdateProfileInput(firstName, lastName string) UpdateProfileInput {
    return UpdateProfileInput{
        FirstName: firstName,
        LastName:  lastName,
    }
}

// MoveTaskInput represents task move data
type MoveTaskInput struct {
    Status   string `json:"status"`
    Position int    `json:"position"`
}

// NewMoveTaskInput creates move task input
func NewMoveTaskInput(status string, position int) MoveTaskInput {
    return MoveTaskInput{
        Status:   status,
        Position: position,
    }
}

// AssignLabelInput represents label assignment data
type AssignLabelInput struct {
    LabelID string `json:"label_id"`
}

// NewAssignLabelInput creates assign label input
func NewAssignLabelInput(labelID string) AssignLabelInput {
    return AssignLabelInput{
        LabelID: labelID,
    }
}
```

## Files to Create

| File | Action |
|------|--------|
| `backend/test/e2e/setup.go` | CREATE |
| `backend/test/e2e/migrations.go` | CREATE |
| `backend/test/e2e/client.go` | CREATE |
| `backend/test/fixtures/testdata.go` | CREATE |

## Verification

```bash
# Verify setup compiles
cd backend && go build ./test/e2e/... ./test/fixtures/...

# Should output nothing on success
```

## Next Steps

After this shared infrastructure is in place, create the individual E2E test files:

1. `11-e2e-auth-flow.md` - Registration and login tests
2. `12-e2e-admin-flow.md` - Complete admin workflow tests
3. `13-e2e-project-flow.md` - Project, task, epic management tests

Each test file will use `SetupTestEnv(t)` to get an isolated database.
