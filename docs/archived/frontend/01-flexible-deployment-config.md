# Frontend: Flexible Deployment Configuration

## Objective

Make the frontend configurable to support three communication patterns with the backend:

1. **Same domain** — `PUBLIC_API_URL=/api`, nginx proxies `/api/*` to backend
2. **Different domains** — `PUBLIC_API_URL=https://api.example.com`, direct cross-origin fetch
3. **API gateway** — `PUBLIC_API_URL=/api` or `PUBLIC_API_URL=https://gateway.example.com/api`

## Problem Analysis

### Current Issues

| File | Issue | Impact |
|------|-------|--------|
| `vite.config.ts:10` | Proxy target hardcoded to `http://localhost:8080` | Dev only works on port 8080 |
| `frontend/.env` | Only has `/api` — no documentation on other patterns | Confusing for deployers |
| `frontend/.env.example` | Same as above | No guidance |
| `frontend/.env.production` | Same `/api` — not wrong, but missing context | No guidance |

### What Already Works

- `PUBLIC_API_URL` in `src/lib/api/real/client.ts` already drives all API calls — this is correct
- `credentials: 'include'` is already set on all requests — cookies work cross-origin when backend allows it
- No code changes needed in API layer for different-domain support — only config changes

---

## Task 1: Make Vite Proxy Target Configurable

**File:** `frontend/vite.config.ts`

Read backend URL from `process.env` (Node.js environment, available at dev-server startup):

```typescript
// Before
target: 'http://localhost:8080',

// After
target: process.env.BACKEND_URL ?? 'http://localhost:8080',
```

Full updated server block:
```typescript
server: {
    proxy: {
        '/api': {
            target: process.env.BACKEND_URL ?? 'http://localhost:8080',
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api/, ''),
            cookieDomainRewrite: 'localhost'
        }
    }
},
```

> The proxy only runs during `npm run dev`. In production, nginx or the gateway handles routing.
> `BACKEND_URL` is a **Node.js env var** (not a `PUBLIC_` Vite var) — it never reaches the browser.

**Verification:** `npm run dev` still works with no `BACKEND_URL` set (falls back to port 8080).

---

## Task 2: Update .env Files

**`frontend/.env`** (local dev):
```env
# API base URL — relative path proxied by Vite to BACKEND_URL
PUBLIC_API_URL=/api
PUBLIC_USE_MOCK_API=false

# Vite dev server proxy target (Node.js only, not exposed to browser)
BACKEND_URL=http://localhost:8080
```

**`frontend/.env.example`**:
```env
# -----------------------------------------------
# PATTERN 1: Same domain (nginx proxies /api/*)
# -----------------------------------------------
PUBLIC_API_URL=/api
PUBLIC_USE_MOCK_API=false

# Dev only: Vite proxy target (ignored in production)
BACKEND_URL=http://localhost:8080

# -----------------------------------------------
# PATTERN 2: Different domains (cross-origin)
# -----------------------------------------------
# PUBLIC_API_URL=https://api.example.com
# (BACKEND_URL not needed — no proxy in prod)

# -----------------------------------------------
# PATTERN 3: API gateway
# -----------------------------------------------
# PUBLIC_API_URL=/api                           # if gateway is same domain
# PUBLIC_API_URL=https://gateway.example.com/api  # if gateway is different domain
```

**`frontend/.env.production`**:
```env
# Set to match your deployment pattern:
#   Same domain or gateway (same host): /api
#   Different domain or gateway (different host): https://api.example.com
PUBLIC_API_URL=/api
PUBLIC_USE_MOCK_API=false
```

---

## Task 3: Document Deployment Patterns

**No code changes** — this is configuration-only guidance baked into `.env.example`.

### How Each Pattern Works

**Pattern 1: Same Domain**
```
https://app.example.com
  /          → frontend (nginx serves static build)
  /api/*     → nginx proxies to backend:8080
```
```env
PUBLIC_API_URL=/api
```
Browser calls `/api/projects` → same origin → nginx routes to backend → no CORS needed.

---

**Pattern 2: Different Domains**
```
https://app.example.com  → frontend
https://api.example.com  → backend
```
```env
PUBLIC_API_URL=https://api.example.com
```
Browser calls `https://api.example.com/projects` → cross-origin → backend must:
- Allow `https://app.example.com` in `CORS_ORIGINS`
- Set `COOKIE_SECURE=true`, `COOKIE_SAMESITE=none`, `COOKIE_DOMAIN=.example.com`

---

**Pattern 3: API Gateway (same domain)**
```
https://app.example.com
  /          → frontend
  /api/*     → API gateway → backend
```
```env
PUBLIC_API_URL=/api
```
Same as Pattern 1 from the frontend's perspective. Gateway handles backend routing.

---

**Pattern 3b: API Gateway (different domain)**
```
https://app.example.com   → frontend
https://gateway.example.com/api → gateway → backend
```
```env
PUBLIC_API_URL=https://gateway.example.com/api
```
Same as Pattern 2 from frontend's perspective. Gateway must handle CORS or forward headers.

---

## Files Changed

| File | Change |
|------|--------|
| `frontend/vite.config.ts` | Read proxy target from `process.env.BACKEND_URL` |
| `frontend/.env` | Add `BACKEND_URL=http://localhost:8080` |
| `frontend/.env.example` | Full documentation of all three patterns |
| `frontend/.env.production` | Updated with guidance comment |

## Success Criteria

- [ ] `npm run dev` works without setting `BACKEND_URL` (defaults to port 8080)
- [ ] `npm run dev` works when `BACKEND_URL=http://localhost:9090` is set
- [ ] `npm run build` succeeds with `PUBLIC_API_URL=/api`
- [ ] `npm run build` succeeds with `PUBLIC_API_URL=https://api.example.com`
- [ ] `.env.example` clearly documents all three deployment patterns

## Deployment Cheat Sheet

| Pattern | `PUBLIC_API_URL` | `BACKEND_URL` (dev only) |
|---------|-----------------|--------------------------|
| Same domain (nginx) | `/api` | `http://localhost:8080` |
| Different domains | `https://api.example.com` | `http://localhost:8080` |
| API gateway (same host) | `/api` | `http://localhost:8080` |
| API gateway (different host) | `https://gateway.example.com/api` | `http://localhost:8080` |

> `BACKEND_URL` is only used by the Vite dev server proxy. It has no effect in production builds.
