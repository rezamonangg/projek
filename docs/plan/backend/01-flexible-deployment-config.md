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

> `CORS_ORIGINS` is a comma-separated string. Parse it into `[]string` after `viper.Unmarshal`.

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
| `internal/common/config.go` | Add `CORSConfig`, `CookieConfig` structs, defaults, env bindings |
| `internal/router/router.go` | Use `cfg.CORS.AllowedOrigins`, `cfg.Auth.SessionTTL`, pass `cfg.Cookie` to handler |
| `internal/auth/handler.go` | Accept `CookieConfig`, use it in both `SetCookie` calls, add `parseSameSite` helper |
| `backend/.env` | Add `CORS_ORIGINS`, `COOKIE_*` vars |
| `backend/.env.example` | Same |
| `backend/.env.production` | Production values |

## Success Criteria

- [ ] `go build ./...` passes with no errors
- [ ] `go test ./...` passes with no regressions
- [ ] `CORS_ORIGINS` env var drives allowed origins (no restart needed beyond env change)
- [ ] `COOKIE_SECURE=true` + `COOKIE_SAMESITE=none` works for cross-domain cookies
- [ ] `COOKIE_DOMAIN=.example.com` enables subdomain cookie sharing
- [ ] Local dev still works with defaults (no env change required)
