# E2E Test Plan: Authentication Flow

**Priority:** HIGH
**Estimated Time:** 1 hour
**Dependencies:** 10-e2e-shared-setup.md

## Test Scope

| Flow | Description |
|------|-------------|
| Register | User registers a new community (becomes admin) |
| Login | User logs in with email/password |
| Logout | User logs out |
| Protected | Access protected endpoints without auth |

## Database Isolation

Each test case creates a **fresh PostgreSQL and Redis container** via `SetupTestEnv(t)`. This ensures:
- No state leakage between tests
- No cleanup needed (containers are destroyed after test)
- True isolation for parallel test execution

## Test File

**File:** `backend/test/e2e/auth_test.go`

```go
package e2e

import (
    "context"
    "net/http"
    "testing"
    "time"

    "github.com/monachy/projek/test/fixtures"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// ============================================================
// TEST SUITE: Registration Flow
// ============================================================

func TestAuth_Register_Success(t *testing.T) {
    // GIVEN: Fresh isolated database
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // WHEN: User registers a new community
    input := fixtures.NewCommunityInput()
    
    resp, err := env.Client.Post("/auth/register", input)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Registration succeeds
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    var result map[string]interface{}
    ParseResponse(t, resp, &result)
    
    // Verify response contains expected fields
    assert.NotEmpty(t, result["community_id"], "Should return community ID")
    assert.NotEmpty(t, result["member_id"], "Should return member ID")
    
    // Store for subsequent tests
    communityID := result["community_id"].(string)
    memberID := result["member_id"].(string)
    
    t.Logf("Created community: %s, member: %s", communityID, memberID)
    
    // Verify database state
    var dbCount int
    err = env.DB.QueryRow(context.Background(),
        "SELECT COUNT(*) FROM communities WHERE id = $1",
        communityID).Scan(&dbCount)
    require.NoError(t, err)
    assert.Equal(t, 1, dbCount, "Community should exist in database")
    
    // Verify admin member created
    err = env.DB.QueryRow(context.Background(),
        "SELECT COUNT(*) FROM members WHERE community_id = $1 AND role = 'admin'",
        communityID).Scan(&dbCount)
    require.NoError(t, err)
    assert.Equal(t, 1, dbCount, "Should have 1 admin member")
}

func TestAuth_Register_DuplicateSlug(t *testing.T) {
    // GIVEN: Database with existing community
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // Create first community
    input := fixtures.NewCommunityInput()
    resp, _ := env.Client.Post("/auth/register", input)
    resp.Body.Close()
    
    // WHEN: Try to register with same slug
    input2 := fixtures.NewCommunityInput()
    input2.AdminEmail = "different@test.com" // Different email
    
    resp2, err := env.Client.Post("/auth/register", input2)
    require.NoError(t, err)
    defer resp2.Body.Close()
    
    // THEN: Should fail with conflict
    assert.Equal(t, http.StatusConflict, resp2.StatusCode)
    
    errResp := ParseError(t, resp2)
    assert.Contains(t, errResp.Error.Message, "slug")
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {
    // GIVEN: Database with existing member
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // Create first community
    input := fixtures.NewCommunityInput()
    resp, _ := env.Client.Post("/auth/register", input)
    resp.Body.Close()
    
    // WHEN: Try to register with same email (even different community)
    input2 := fixtures.NewCommunityInput()
    input2.Slug = "different-slug"
    
    resp2, err := env.Client.Post("/auth/register", input2)
    require.NoError(t, err)
    defer resp2.Body.Close()
    
    // THEN: Should fail
    assert.NotEqual(t, http.StatusCreated, resp2.StatusCode)
}

func TestAuth_Register_ValidationErrors(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    tests := []struct {
        name       string
        input      fixtures.RegisterCommunityInput
        expectCode int
    }{
        {
            name: "missing email",
            input: fixtures.RegisterCommunityInput{
                Name:          "Test",
                Slug:          "test",
                AdminPassword: "password123",
                FirstName:     "Admin",
                LastName:      "User",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "invalid email",
            input: fixtures.RegisterCommunityInput{
                Name:          "Test",
                Slug:          "test",
                AdminEmail:    "not-an-email",
                AdminPassword: "password123",
                FirstName:     "Admin",
                LastName:      "User",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "missing password",
            input: fixtures.RegisterCommunityInput{
                Name:       "Test",
                Slug:       "test",
                AdminEmail: "admin@test.com",
                FirstName:  "Admin",
                LastName:   "User",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "short password",
            input: fixtures.RegisterCommunityInput{
                Name:          "Test",
                Slug:          "test",
                AdminEmail:    "admin@test.com",
                AdminPassword: "short",
                FirstName:     "Admin",
                LastName:      "User",
            },
            expectCode: http.StatusBadRequest,
        },
        {
            name: "missing name",
            input: fixtures.RegisterCommunityInput{
                Slug:          "test",
                AdminEmail:    "admin@test.com",
                AdminPassword: "password123",
                FirstName:     "Admin",
                LastName:      "User",
            },
            expectCode: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := env.Client.Post("/auth/register", tt.input)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, tt.expectCode, resp.StatusCode)
        })
    }
}

// ============================================================
// TEST SUITE: Login Flow
// ============================================================

func TestAuth_Login_Success(t *testing.T) {
    // GIVEN: Registered community with admin
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerInput := fixtures.NewCommunityInput()
    resp, _ := env.Client.Post("/auth/register", registerInput)
    resp.Body.Close()
    
    // WHEN: Admin logs in
    loginInput := fixtures.NewLoginInput("admin@test.com", "password123")
    
    resp, err := env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Login succeeds with session
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    var result map[string]interface{}
    ParseResponse(t, resp, &result)
    
    assert.NotEmpty(t, result["session_id"], "Should return session ID")
    
    member := result["member"].(map[string]interface{})
    assert.Equal(t, "admin@test.com", member["email"])
    assert.Equal(t, "admin", member["role"])
    assert.Equal(t, "Admin", member["first_name"])
    assert.Equal(t, "User", member["last_name"])
    
    // Verify session cookie is set
    sessionCookie := env.Client.GetSessionCookie()
    assert.NotNil(t, sessionCookie, "Session cookie should be set")
    assert.NotEmpty(t, sessionCookie.Value, "Session cookie should have value")
    
    t.Logf("Login successful, session: %s", result["session_id"])
}

func TestAuth_Login_WrongPassword(t *testing.T) {
    // GIVEN: Registered community
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerInput := fixtures.NewCommunityInput()
    resp, _ := env.Client.Post("/auth/register", registerInput)
    resp.Body.Close()
    
    // WHEN: Try to login with wrong password
    loginInput := fixtures.NewLoginInput("admin@test.com", "wrongpassword")
    
    resp, err := env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should be unauthorized
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    
    // No session cookie should be set
    sessionCookie := env.Client.GetSessionCookie()
    assert.Nil(t, sessionCookie, "No session cookie on failed login")
}

func TestAuth_Login_NonExistentUser(t *testing.T) {
    // GIVEN: Empty database
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // WHEN: Try to login with non-existent user
    loginInput := fixtures.NewLoginInput("nobody@test.com", "password123")
    
    resp, err := env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should be unauthorized
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_Login_ValidationErrors(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    tests := []struct {
        name       string
        input      fixtures.LoginInput
        expectCode int
    }{
        {
            name:       "missing email",
            input:      fixtures.NewLoginInput("", "password"),
            expectCode: http.StatusBadRequest,
        },
        {
            name:       "missing password",
            input:      fixtures.NewLoginInput("admin@test.com", ""),
            expectCode: http.StatusBadRequest,
        },
        {
            name:       "invalid email format",
            input:      fixtures.NewLoginInput("not-an-email", "password"),
            expectCode: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := env.Client.Post("/auth/login", tt.input)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, tt.expectCode, resp.StatusCode)
        })
    }
}

// ============================================================
// TEST SUITE: Logout Flow
// ============================================================

func TestAuth_Logout_Success(t *testing.T) {
    // GIVEN: Logged in admin
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerAndLogin(t, env)
    
    // Verify session exists in Redis
    sessionCookie := env.Client.GetSessionCookie()
    require.NotNil(t, sessionCookie)
    
    // WHEN: Admin logs out
    resp, err := env.Client.Post("/auth/logout", nil)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Logout succeeds
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Verify session is removed from Redis
    ctx := context.Background()
    exists := env.Redis.Exists(ctx, "session:"+sessionCookie.Value).Val()
    assert.Equal(t, int64(0), exists, "Session should be removed from Redis")
}

func TestAuth_Logout_WithoutSession(t *testing.T) {
    // GIVEN: Fresh client without session
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // WHEN: Try to logout without being logged in
    resp, err := env.Client.Post("/auth/logout", nil)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should still succeed (idempotent)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// ============================================================
// TEST SUITE: Protected Endpoints
// ============================================================

func TestAuth_ProtectedEndpoint_WithoutAuth(t *testing.T) {
    // GIVEN: Fresh client without session
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // WHEN: Access protected endpoint without auth
    resp, err := env.Client.Get("/auth/me")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should be unauthorized
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ProtectedEndpoint_WithAuth(t *testing.T) {
    // GIVEN: Logged in admin
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerAndLogin(t, env)
    
    // WHEN: Access protected endpoint with auth
    resp, err := env.Client.Get("/auth/me")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should succeed
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    var result map[string]interface{}
    ParseResponse(t, resp, &result)
    
    assert.Equal(t, "admin@test.com", result["email"])
    assert.Equal(t, "admin", result["role"])
}

func TestAuth_ProtectedEndpoint_AfterLogout(t *testing.T) {
    // GIVEN: Admin who logged out
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerAndLogin(t, env)
    
    // Logout
    resp, _ := env.Client.Post("/auth/logout", nil)
    resp.Body.Close()
    
    // WHEN: Try to access protected endpoint
    resp, err := env.Client.Get("/auth/me")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should be unauthorized
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ProtectedEndpoint_ExpiredSession(t *testing.T) {
    // GIVEN: Admin with session
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    registerAndLogin(t, env)
    
    sessionCookie := env.Client.GetSessionCookie()
    require.NotNil(t, sessionCookie)
    
    // Simulate session expiration by deleting from Redis
    ctx := context.Background()
    env.Redis.Del(ctx, "session:"+sessionCookie.Value)
    
    // WHEN: Try to access protected endpoint
    resp, err := env.Client.Get("/auth/me")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // THEN: Should be unauthorized
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// ============================================================
// TEST SUITE: Complete Auth Flow
// ============================================================

func TestAuth_CompleteFlow(t *testing.T) {
    // GIVEN: Fresh isolated database
    env := SetupTestEnv(t)
    defer env.Cleanup()
    
    // STEP 1: Register
    t.Log("Step 1: Register community")
    registerInput := fixtures.NewCommunityInput()
    
    resp, err := env.Client.Post("/auth/register", registerInput)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, resp.StatusCode)
    
    var registerResult map[string]interface{}
    ParseResponse(t, resp, &registerResult)
    
    communityID := registerResult["community_id"].(string)
    env.SetCommunityID(communityID)
    t.Logf("Community created: %s", communityID)
    
    // STEP 2: Login
    t.Log("Step 2: Login as admin")
    loginInput := fixtures.NewLoginInput("admin@test.com", "password123")
    
    resp, err = env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    t.Log("Login successful")
    
    // STEP 3: Access protected endpoint
    t.Log("Step 3: Access protected endpoint")
    resp, err = env.Client.Get("/auth/me")
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    var meResult map[string]interface{}
    ParseResponse(t, resp, &meResult)
    
    assert.Equal(t, "admin@test.com", meResult["email"])
    assert.Equal(t, "admin", meResult["role"])
    t.Log("Protected endpoint accessible")
    
    // STEP 4: Logout
    t.Log("Step 4: Logout")
    resp, err = env.Client.Post("/auth/logout", nil)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    t.Log("Logout successful")
    
    // STEP 5: Verify cannot access protected endpoint
    t.Log("Step 5: Verify protected endpoint is now inaccessible")
    resp, err = env.Client.Get("/auth/me")
    require.NoError(t, err)
    require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    
    t.Log("Complete auth flow verified")
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func registerAndLogin(t *testing.T, env *TestEnv) {
    // Register
    registerInput := fixtures.NewCommunityInput()
    resp, err := env.Client.Post("/auth/register", registerInput)
    require.NoError(t, err)
    resp.Body.Close()
    
    // Login
    loginInput := fixtures.NewLoginInput("admin@test.com", "password123")
    resp, err = env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    resp.Body.Close()
    
    // Verify logged in
    sessionCookie := env.Client.GetSessionCookie()
    require.NotNil(t, sessionCookie, "Should have session cookie after login")
}
```

## Test Cases Summary

| Test | DB Isolation | Description |
|------|--------------|-------------|
| `TestAuth_Register_Success` | ✅ Fresh DB | Register creates community + admin |
| `TestAuth_Register_DuplicateSlug` | ✅ Fresh DB | Duplicate slug returns 409 |
| `TestAuth_Register_DuplicateEmail` | ✅ Fresh DB | Duplicate email returns error |
| `TestAuth_Register_ValidationErrors` | ✅ Fresh DB | Invalid inputs return 400 |
| `TestAuth_Login_Success` | ✅ Fresh DB | Valid credentials return session |
| `TestAuth_Login_WrongPassword` | ✅ Fresh DB | Wrong password returns 401 |
| `TestAuth_Login_NonExistentUser` | ✅ Fresh DB | Unknown user returns 401 |
| `TestAuth_Logout_Success` | ✅ Fresh DB | Logout clears session |
| `TestAuth_ProtectedEndpoint_*` | ✅ Fresh DB | Auth-required endpoints |
| `TestAuth_CompleteFlow` | ✅ Fresh DB | Full register→login→logout flow |

## Running Tests

```bash
# Run auth E2E tests
cd backend && go test -v -tags=e2e -run TestAuth ./test/e2e/...

# Run single test
cd backend && go test -v -tags=e2e -run TestAuth_CompleteFlow ./test/e2e/...

# Run with timeout (containers take time)
cd backend && go test -v -tags=e2e -timeout 10m -run TestAuth ./test/e2e/...
```

## Expected Output

```
=== RUN   TestAuth_Register_Success
    auth_test.go:45: Created community: abc123..., member: def456...
--- PASS: TestAuth_Register_Success (5.23s)
=== RUN   TestAuth_Login_Success
    auth_test.go:112: Login successful, session: xyz789...
--- PASS: TestAuth_Login_Success (4.87s)
=== RUN   TestAuth_CompleteFlow
    auth_test.go:312: Step 1: Register community
    auth_test.go:321: Community created: abc123...
    auth_test.go:324: Step 2: Login as admin
    auth_test.go:332: Login successful
    auth_test.go:335: Step 3: Access protected endpoint
    auth_test.go:345: Protected endpoint accessible
    auth_test.go:348: Step 4: Logout
    auth_test.go:354: Logout successful
    auth_test.go:358: Step 5: Verify protected endpoint is now inaccessible
    auth_test.go:363: Complete auth flow verified
--- PASS: TestAuth_CompleteFlow (8.45s)
PASS
```

## Files to Create

| File | Action |
|------|--------|
| `backend/test/e2e/auth_test.go` | CREATE |
