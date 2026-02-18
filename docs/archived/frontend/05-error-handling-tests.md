# Frontend Testing: Error Handling and Loading States

## Overview

E2E tests for error handling, loading states, and network failure scenarios across all pages.

## Task 1: Global Error Handling

**Test Scenarios:**
- TC-01: Toast notifications for errors
  - Trigger API error (e.g., invalid project ID)
  - Verify toast notification appears
  - Verify error message is user-friendly
  - Verify toast auto-dismisses or can be closed

- TC-02: Error boundary catches crashes
  - If possible, trigger component error
  - Verify error boundary shows fallback UI
  - Verify error logged correctly

- TC-03: Offline indicator
  - Go offline (network throttling)
  - Verify offline indicator shown
  - Come back online
  - Verify indicator hidden

**Coverage:**
- Toast system
- Error boundary component
- Offline detection

## Task 2: Loading States

**Test Scenarios:**
- TC-01: Dashboard loading state
  - Login via API
  - Navigate to `/dashboard`
  - Verify "Loading..." text appears initially
  - Verify content replaces loader when data loaded

- TC-02: Projects skeleton loading
  - Login via API
  - Navigate to `/projects`
  - Verify `SkeletonCard` components shown
  - Verify skeletons replaced with actual content
  - Verify smooth transition

- TC-03: Button loading states
  - Click submit button on any form
  - Verify button shows loading spinner
  - Verify button disabled during loading
  - Verify button returns to normal after completion

**Coverage:**
- Loading indicators
- Skeleton components
- Button loading states
- Smooth transitions

## Task 3: Network Error Scenarios

**Test Scenarios:**
- TC-01: API timeout handling
  - Throttle network to very slow
  - Attempt operation
  - Verify loading state persists
  - Verify timeout error handled gracefully

- TC-02: 500 Server Error
  - Mock or trigger 500 error
  - Verify user-friendly error message
  - Verify retry option if available

- TC-03: 404 Not Found
  - Navigate to non-existent resource
  - Verify 404 page or error message
  - Verify navigation options provided

- TC-04: 403 Forbidden
  - Attempt to access resource without permission
  - Verify 403 error handled
  - Verify appropriate message shown

**Coverage:**
- HTTP error codes
- Error message display
- Recovery options

## Task 4: Form Validation Errors

**Test Scenarios:**
- TC-01: Login form validation display
  - Submit empty login form
  - Verify inline validation errors
  - Verify errors clear when field corrected

- TC-02: Create project validation
  - Submit empty project name
  - Verify error message below field
  - Verify field highlighted in error state

- TC-03: Server-side validation errors
  - Trigger server validation error (e.g., duplicate email)
  - Verify error displayed inline or in toast

**Coverage:**
- Client-side validation
- Server-side validation
- Error display patterns

## Task 5: Session Expiration Flow

**Test Scenarios:**
- TC-01: Session expires during usage
  - Login via API
  - Navigate to protected page
  - Invalidate session (clear cookies)
  - Attempt API operation
  - Verify redirect to login
  - Verify error message about session expiration

- TC-02: Token refresh (if implemented)
  - Test automatic token refresh
  - Verify seamless continuation

**Coverage:**
- 401 handling
- Session management
- Redirect flow

## File Changes

```
frontend/tests/e2e/
├── error-handling.spec.ts (CREATE NEW)
└── loading-states.spec.ts (CREATE NEW - optional, can combine)
```

## Test Utilities

```typescript
// Network throttling helpers
async function goOffline(page) {
  await page.context().setOffline(true);
}

async function goOnline(page) {
  await page.context().setOffline(false);
}

async function throttleNetwork(page, speed: 'slow' | 'fast') {
  // Use Playwright's network throttling
}

// Error mocking
async function mockAPIError(page, endpoint: string, status: number) {
  await page.route(endpoint, route => {
    route.fulfill({ status });
  });
}
```

## Expected Test Count

- Global error handling: 3
- Loading states: 3
- Network errors: 4
- Form validation: 3
- Session expiration: 2
- **Total: 15 tests**

## Notes

- These tests apply patterns across all pages
- Can be combined with specific page tests if preferred
- Network throttling requires Playwright advanced features
