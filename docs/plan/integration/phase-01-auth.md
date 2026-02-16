# Integration Phase 1: Auth APIs

Replace mock auth APIs with real backend endpoints.

## Prerequisites
- Backend auth endpoints implemented and running
- Backend CORS configured for frontend origin

## Task 1.1: Real Auth API Implementation

**Commits:**
- `feat: add real auth api client`
- `feat: add axios/fetch http client`
- `feat: add request/response interceptors`
- `feat: add error handling`

**Endpoints:**
- POST /api/auth/login
- POST /api/auth/logout
- GET /api/auth/me

## Task 1.2: Session Management

**Commits:**
- `feat: add session cookie handling`
- `feat: update auth store for real sessions`
- `feat: add session expiration handling`

## Task 1.3: Switch to Real Auth

**Commits:**
- `feat: switch auth api to real implementation`
- `test: verify login flow with backend`
- `test: verify logout flow`
- `test: verify session persistence`

## Task 1.4: Remove Mock Auth

**Commits:**
- `chore: remove mock auth api`
- `chore: update documentation`

---

## Total Integration Commits (Phase 1): 10
