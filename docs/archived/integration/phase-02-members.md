# Integration Phase 2: Member APIs

Replace mock member APIs with real backend endpoints.

## Prerequisites
- Backend member endpoints implemented
- Auth integration complete

## Task 2.1: Real Member API Implementation

**Commits:**
- `feat: add real member api client`
- `feat: add member list endpoint integration`
- `feat: add invite member endpoint integration`
- `feat: add profile update endpoint integration`

**Endpoints:**
- GET /api/members
- POST /api/members/invite
- GET /api/members/:id
- PUT /api/members/:id
- PUT /api/members/me

## Task 2.2: Switch to Real Members

**Commits:**
- `feat: switch member api to real implementation`
- `test: verify member list loading`
- `test: verify invite member flow`
- `test: verify profile update`

## Task 2.3: Remove Mock Members

**Commits:**
- `chore: remove mock member api`

---

## Total Integration Commits (Phase 2): 7
