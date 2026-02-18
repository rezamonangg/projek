# Backend: Flexible Deployment Configuration

## Objective

Make the backend configurable to support three frontend-to-backend communication patterns:

1. **Same domain** — nginx proxies `/api/*` to backend on the same host
2. **Different domains** — frontend and backend on separate origins (CORS + cross-site cookies)
3. **API gateway** — gateway sits between frontend and backend, same or different domain

## Problem Analysis

### Current Hardcoded Values

| File | Line | What | Problem |
|------|------|------|---------|
| `internal/router/router.go` | 33 | CORS origins: `localhost:3000, :5173, :8080` | Breaks any production domain |
| `internal/auth/handler.go` | 64 | `Secure: false` | Breaks HTTPS deployments |
| `internal/auth/handler.go` | 65 | `SameSite: Lax` | Breaks cross-site cookie (different domain) |
| `internal/auth/handler.go` | 59-67 | No `Domain` on cookie | Cannot share cookie across subdomains |
| `internal/auth/handler.go` | 66 | `MaxAge: 7*24*60*60` hardcoded | Duplicate of config value |
| `internal/router/router.go` | 46 | Session TTL `7*24*time.Hour` hardcoded | Ignores `cfg.Auth.SessionTTL` |

---

## Task 1: Extend Config Struct

**File:** `internal/common/config.go`

Add two new config structs:

```go
type CORSConfig struct {
    AllowedOrigins []string // e.g. ["https://app.example.com"]
}

type CookieConfig struct {
    Secure   bool   // true in production (HTTPS)
    Domain   string // empty = exact host; ".example.com" = subdomain sharing
    SameSite string // "lax", "strict", "none"
}
```

Add to `Config`:
```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Auth     AuthConfig
    Email    EmailConfig
    CORS     CORSConfig    // NEW
    Cookie   CookieConfig  // NEW
}
```

Default values:
```go
CORS: CORSConfig{
    AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
},
Cookie: CookieConfig{
    Secure:   false,
    Domain:   "",
    SameSite: "lax",
},
```

Env var bindings (via viper):
```go
viper.RegisterAlias("cors.allowed_origins", "CORS_ORIGINS")     // comma-separated
viper.RegisterAlias("cookie.secure",         "COOKIE_SECURE")
viper.RegisterAlias("cookie.domain",         "COOKIE_DOMAIN")
viper.RegisterAlias("cookie.same_site",      "COOKIE_SAMESITE")
```

Parse comma-separated CORS_ORIGINS:
```go
if origins := viper.GetString("cors.allowed_origins"); origins != "" {
    cfg.CORS.AllowedOrigins = strings.Split(origins, ",")
}
```

Add validation for SameSite=none requiring Secure=true:
```go
if cfg.Cookie.SameSite == "none" && !cfg.Cookie.Secure {
    log.Fatal().Msg("SameSite=none requires Secure=true")
}
```

**Verification:** `go build ./...` passes with new struct fields.

---

## Task 2: Use Config in CORS Middleware

**File:** `internal/router/router.go`

Replace hardcoded origins:

```go
// Before
AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"},

// After
AllowedOrigins: cfg.CORS.AllowedOrigins,
```

**Verification:** `go build ./...` passes.

---

## Task 3: Use Config in Session Store

**File:** `internal/router/router.go`

Replace hardcoded TTL:

```go
// Before
sessionStore := auth.NewSessionStore(redis, 7*24*time.Hour)

// After
sessionStore := auth.NewSessionStore(redis, cfg.Auth.SessionTTL)
```

**Verification:** `go build ./...` passes.

---

## Task 4: Cookie Config in Auth Handler

**File:** `internal/auth/handler.go`

Pass `CookieConfig` into the handler:

```go
type Handler struct {
    authService   *AuthService
    memberService *member.Service
    cookieCfg     common.CookieConfig  // NEW
}

func NewHandler(authService *AuthService, memberService *member.Service, cookieCfg common.CookieConfig) *Handler {
    return &Handler{
        authService:   authService,
        memberService: memberService,
        cookieCfg:     cookieCfg,
    }
}
```

Add a helper to resolve `http.SameSite` from string:

```go
func parseSameSite(s string) http.SameSite {
    switch strings.ToLower(s) {
    case "strict":
        return http.SameSiteStrictMode
    case "none":
        return http.SameSiteNoneMode
    default:
        return http.SameSiteLaxMode
    }
}
```

Use in `Login`:
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    session.ID,
    Path:     "/",
    HttpOnly: true,
    Secure:   h.cookieCfg.Secure,
    Domain:   h.cookieCfg.Domain,
    SameSite: parseSameSite(h.cookieCfg.SameSite),
    MaxAge:   int(h.authService.sessionTTL.Seconds()),
})
```

Use in `Logout` (same attributes, `MaxAge: -1`):
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    "",
    Path:     "/",
    HttpOnly: true,
    Secure:   h.cookieCfg.Secure,
    Domain:   h.cookieCfg.Domain,
    SameSite: parseSameSite(h.cookieCfg.SameSite),
    MaxAge:   -1,
})
```

Expose `sessionTTL` in `AuthService` (or pass TTL separately to `NewHandler`).

**Verification:** `go build ./...` passes, `go test ./internal/auth/...` passes.

---

## Task 5: Update NewHandler call in router.go

**File:** `internal/router/router.go`

```go
authHandler := auth.NewHandler(authService, memberService, cfg.Cookie)
```

**Verification:** `go build ./...` passes.

---

## Task 6: Update .env Files

**`backend/.env` and `backend/.env.example`** — add:
```env
# CORS
CORS_ORIGINS=http://localhost:5173,http://localhost:3000

# Cookie
COOKIE_SECURE=false
COOKIE_DOMAIN=
COOKIE_SAMESITE=lax
```

**`backend/.env.production`** — add:
```env
# CORS — set to your frontend origin(s)
CORS_ORIGINS=https://app.example.com

# Cookie — always secure in production
COOKIE_SECURE=true
COOKIE_DOMAIN=
COOKIE_SAMESITE=lax
```

---

## Task 7: Update E2E Test Setup

**File:** `test/e2e/setup.go`

Update the Config initialization to include CORS and Cookie configs:

```go
cfg := &common.Config{
    Server: common.ServerConfig{
        Port:         8080,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    },
    Database: common.DatabaseConfig{
        Host:     "localhost",
        Port:     5432,
        User:     "test",
        Password: "test",
        Database: dbName,
        PoolSize: 5,
    },
    Redis: common.RedisConfig{
        Host: redisHost,
        Port: 6379,
    },
    Auth: common.AuthConfig{
        SessionTTL:     24 * time.Hour,
        PasswordMinLen: 8,
    },
    Email: common.EmailConfig{
        SMTPHost:  "localhost",
        SMTPPort:  587,
        FromEmail: "test@localhost",
        FromName:  "Test",
    },
    CORS: common.CORSConfig{
        AllowedOrigins: []string{"http://localhost:5173"},
    },
    Cookie: common.CookieConfig{
        Secure:   false,
        Domain:   "",
        SameSite: "lax",
    },
}
```

**Verification:** `go test -tags=e2e ./test/...` passes.

---

## Task 8: Update E2E Auth Tests

**File:** `test/e2e/auth_test.go`

Add test to verify cookie attributes are set correctly:

```go
func TestAuth_Login_SetsCorrectCookieAttributes(t *testing.T) {
    env := SetupTestEnv(t)
    defer env.Cleanup()

    ctx := context.Background()

    err := env.ExecDB(ctx, `
        INSERT INTO communities (id, name, slug, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
    `, fixtures.CommunityID, "Test Community", "test-community")
    require.NoError(t, err)

    passwordHash := "$2a$10$sKr8GxHBJAbOdab9Ma.6NOVV9XFORV6bKg1VOpH1guE7rlv4SucO."
    err = env.ExecDB(ctx, `
        INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
    `, fixtures.MemberID, fixtures.CommunityID, "admin@test.com", passwordHash, "Admin", "User", "admin", true)
    require.NoError(t, err)

    loginInput := fixtures.E2ENewLoginInput("admin@test.com", "password123")

    resp, err := env.Client.Post("/auth/login", loginInput)
    require.NoError(t, err)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    sessionCookie := env.Client.GetSessionCookie()
    require.NotNil(t, sessionCookie, "Session cookie should be set after login")
    
    // Verify cookie attributes match test config (Secure: false, SameSite: lax)
    assert.Equal(t, "session_id", sessionCookie.Name)
    assert.NotEmpty(t, sessionCookie.Value, "Session cookie should have a value")
    assert.False(t, sessionCookie.Secure, "Test config uses Secure=false")
    assert.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite, "Test config uses SameSite=lax")
}
```

**Verification:** `go test -tags=e2e ./test/e2e/...` passes.

---

## Task 9: (Optional) Unit Tests for Auth Handler

**File:** `internal/auth/handler_test.go` (create new file)

Add unit tests for:

```go
func TestHandler_Login_SetsCookieWithConfig(t *testing.T) {
    // Test that Login sets cookie with correct Secure, Domain, SameSite from CookieConfig
}

func TestHandler_Logout_ClearsCookieWithSameAttributes(t *testing.T) {
    // Test that Logout uses same Secure/Domain/SameSite as Login
}

func TestParseSameSite(t *testing.T) {
    tests := []struct {
        input    string
        expected http.SameSite
    }{
        {"lax", http.SameSiteLaxMode},
        {"Lax", http.SameSiteLaxMode},
        {"strict", http.SameSiteStrictMode},
        {"Strict", http.SameSiteStrictMode},
        {"none", http.SameSiteNoneMode},
        {"None", http.SameSiteNoneMode},
        {"invalid", http.SameSiteLaxMode}, // default
    }
    
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := parseSameSite(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Verification:** `go test ./internal/auth/...` passes.

---

## Deployment Cheat Sheet

| Pattern | `CORS_ORIGINS` | `COOKIE_SECURE` | `COOKIE_DOMAIN` | `COOKIE_SAMESITE` |
|---------|---------------|----------------|----------------|------------------|
| Same domain (nginx) | `https://app.example.com` | `true` | _(empty)_ | `lax` |
| Different domains | `https://app.example.com` | `true` | `.example.com` | `none` |
| API gateway (same domain) | `https://app.example.com` | `true` | _(empty)_ | `lax` |
| Local dev | `http://localhost:5173` | `false` | _(empty)_ | `lax` |

> **Note:** `SameSite=None` requires `Secure=true`. Browsers reject `SameSite=None; Secure=false`.

---

## Files Changed

| File | Change |
|------|--------|
| `internal/common/config.go` | Add `CORSConfig`, `CookieConfig` structs, defaults, env bindings, validation |
| `internal/router/router.go` | Use `cfg.CORS.AllowedOrigins`, `cfg.Auth.SessionTTL`, pass `cfg.Cookie` to handler |
| `internal/auth/handler.go` | Accept `CookieConfig`, use it in both `SetCookie` calls, add `parseSameSite` helper |
| `backend/.env` | Add `CORS_ORIGINS`, `COOKIE_*` vars |
| `backend/.env.example` | Same |
| `backend/.env.production` | Production values |
| `test/e2e/setup.go` | Add `CORS` and `Cookie` to test Config struct |
| `test/e2e/auth_test.go` | Add `TestAuth_Login_SetsCorrectCookieAttributes` test |
| `internal/auth/handler_test.go` | (Optional) Unit tests for cookie handling |

## Success Criteria

- [ ] `go build ./...` passes with no errors
- [ ] `go test ./...` passes with no regressions
- [ ] `go test -tags=e2e ./test/...` passes with no regressions
- [ ] `CORS_ORIGINS` env var drives allowed origins (no restart needed beyond env change)
- [ ] `COOKIE_SECURE=true` + `COOKIE_SAMESITE=none` works for cross-domain cookies
- [ ] `COOKIE_DOMAIN=.example.com` enables subdomain cookie sharing
- [ ] SameSite=none without Secure=true fails validation with clear error
- [ ] Local dev still works with defaults (no env change required)
