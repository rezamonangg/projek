# Frontend Testing: Authentication Flow

## Overview

Comprehensive E2E tests for authentication flows with real backend integration.

## Task 1: Login Flow Tests

**Test Scenarios:**
- TC-01: Successful login with valid credentials
  - Navigate to `/login`
  - Enter valid email and password
  - Submit form
  - Verify redirect to `/dashboard`
  - Verify user is authenticated

- TC-02: Login with invalid credentials
  - Navigate to `/login`
  - Enter invalid password
  - Submit form
  - Verify error message displayed
  - Verify stays on login page

- TC-03: Login with invalid email format
  - Navigate to `/login`
  - Enter malformed email
  - Submit form
  - Verify email validation error
  - Verify form not submitted

- TC-04: Login with empty fields
  - Navigate to `/login`
  - Submit empty form
  - Verify required field errors
  - Verify API not called

**API Coverage:**
- `POST /auth/login` - Valid credentials
- `POST /auth/login` - Invalid credentials (401)
- Form validation (client-side)

## Task 2: Logout Flow Tests

**Test Scenarios:**
- TC-01: Successful logout
  - Login first (via API)
  - Navigate to protected page
  - Click logout
  - Verify redirect to `/login`
  - Verify cannot access protected routes

- TC-02: Session expiration
  - Login with short-lived session (if possible)
  - Wait for session to expire
  - Try to access protected page
  - Verify redirect to login
  - Verify error message shown

**API Coverage:**
- `POST /auth/logout`
- `GET /auth/me` - Returns 401 when logged out

## Task 3: Session Management Tests

**Test Scenarios:**
- TC-01: Session persistence across page reload
  - Login successfully
  - Reload page
  - Verify still authenticated
  - Verify user data displayed

- TC-02: Protected route access when authenticated
  - Login via API
  - Navigate to `/dashboard`
  - Verify dashboard loads
  - Verify no redirect to login

- TC-03: Protected route access when NOT authenticated
  - Clear cookies/localStorage
  - Navigate to `/projects`
  - Verify redirect to `/login`

**API Coverage:**
- `GET /auth/me` - Returns user when authenticated
- Cookie/session handling

## Task 4: Registration Flow Tests

**Test Scenarios:**
- TC-01: Registration form validation
  - Navigate to `/register`
  - Submit empty form
  - Verify all required field errors

- TC-02: Password confirmation validation
  - Enter password and different confirm password
  - Verify mismatch error

- TC-03: Community slug validation
  - Enter invalid slug format
  - Verify validation error

**Note:** Registration is not fully implemented in backend, tests should reflect current state

## File Changes

```
frontend/tests/e2e/
├── auth.spec.ts (UPDATE - Add comprehensive tests)
```

## Test Setup Requirements

```typescript
// Test setup utilities to add
async function loginViaAPI(page, email, password) {
  // Direct API call to login endpoint
  // Set cookies/session
}

async function clearSession(page) {
  // Clear cookies
  // Clear localStorage
}
```

## Expected Test Count

- Login tests: 4
- Logout tests: 2
- Session tests: 3
- Registration tests: 3
- **Total: 12 new tests**
