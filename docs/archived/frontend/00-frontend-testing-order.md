# Frontend Testing Implementation Order

## Summary

This plan addresses the gap between frontend tests and the recently implemented backend integration. Current tests only cover basic unauthenticated redirects. These plans add comprehensive E2E tests for real API interactions.

## Implementation Priority

| Order | Plan | Priority | Est. Time | Tests |
|-------|------|----------|-----------|-------|
| 1 | 01-auth-tests.md | **P0** | 2 hours | 12 |
| 2 | 03-projects-tests.md | **P0** | 2 hours | 11 |
| 3 | 02-dashboard-tests.md | **P1** | 1.5 hours | 8 |
| 4 | 04-kanban-tests.md | **P1** | 2 hours | 12 |
| 5 | 05-error-handling-tests.md | **P2** | 1.5 hours | 15 |
| **Total** | | | **9 hours** | **58** |

## Phase 1: Critical Path (Must Have)

### Step 1: Authentication Tests (01-auth-tests.md)

**Why first:** All other tests depend on authentication

**Key deliverables:**
- Login flow with real API
- Logout functionality
- Session management
- Protected route access

**Test setup needed:**
```typescript
// helpers/auth.ts
async function loginViaAPI(page, email, password)
async function clearSession(page)
```

### Step 2: Projects Tests (03-projects-tests.md)

**Why second:** Core functionality, most user interactions

**Key deliverables:**
- Projects list loading from API
- Create project form
- Project card interactions
- Replace misleading existing tests

**APIs covered:**
- `GET /projects`
- `POST /projects`

## Phase 2: Key Features (Should Have)

### Step 3: Dashboard Tests (02-dashboard-tests.md)

**Key deliverables:**
- Stats loading from API
- Recent projects display
- Task distribution charts
- Loading states

**APIs covered:**
- `GET /admin/dashboard/stats`
- `GET /projects` (with limit)

### Step 4: Kanban Tests (04-kanban-tests.md)

**Key deliverables:**
- Board loading with tasks
- Task cards display
- Drag and drop (if implemented)
- Create task from board

**APIs covered:**
- `GET /boards/{id}`
- `GET /boards/{boardId}/tasks`
- `POST /boards/{boardId}/tasks`

## Phase 3: Edge Cases (Nice to Have)

### Step 5: Error Handling Tests (05-error-handling-tests.md)

**Key deliverables:**
- Toast notifications
- Loading states
- Network error scenarios
- Form validation errors
- Session expiration

**Coverage:**
- Error boundaries
- HTTP error codes (404, 403, 500)
- Offline detection
- Validation patterns

## File Changes Summary

### New Files
```
frontend/tests/e2e/
├── dashboard.spec.ts (NEW)
├── error-handling.spec.ts (NEW)
└── helpers/
    ├── auth.ts (NEW - test utilities)
    └── api.ts (NEW - API helpers)
```

### Updated Files
```
frontend/tests/e2e/
├── auth.spec.ts (UPDATE - comprehensive tests)
├── projects.spec.ts (UPDATE - replace misleading tests)
└── kanban.spec.ts (UPDATE - real board tests)
```

## Test Data Strategy

### Setup (beforeEach)
1. Create test user via API (or use demo credentials)
2. Login via API (set cookies/session)
3. Create test data (projects, boards, tasks)

### Teardown (afterEach)
1. Delete test data via API
2. Clear session/cookies

### Example
```typescript
test.beforeEach(async ({ page }) => {
  // Login via API
  await loginViaAPI(page, 'admin@example.com', 'admin123');
  
  // Create test project
  await createTestProject(page, { name: 'Test Project' });
});

test.afterEach(async ({ page }) => {
  // Cleanup
  await cleanupTestProjects(page);
  await clearSession(page);
});
```

## Testing Approach

### Real API vs Mock
- **Use real API** for integration testing
- Tests run against actual backend
- Requires backend to be running
- More reliable but slower

### Authentication Strategy
- Use `loginViaAPI()` helper to authenticate
- Set cookies directly instead of going through UI
- Faster tests, more reliable

### Parallel Execution
- Playwright runs tests in parallel
- Each test should be isolated
- Don't share state between tests
- Use unique test data (timestamps, UUIDs)

## Success Criteria

- [ ] All new tests pass
- [ ] Existing tests still pass (if kept)
- [ ] Tests cover all API endpoints used by frontend
- [ ] Tests verify loading states
- [ ] Tests verify error handling
- [ ] Test utilities are reusable
- [ ] Documentation updated

## Commit Strategy

Recommended commit sequence:

```
1. test: add authentication test utilities and helpers
2. test: add comprehensive auth flow tests
3. test: add projects page tests with API integration
4. test: add dashboard tests with stats loading
5. test: add kanban board tests
6. test: add error handling and loading state tests
7. test: remove or update obsolete tests
```

## Risks and Mitigation

| Risk | Mitigation |
|------|------------|
| Backend not available | Ensure backend running before tests; use health check |
| Test data conflicts | Use unique identifiers; clean up after tests |
| Slow test execution | Parallel execution; API login instead of UI |
| Flaky tests | Add retries; wait for network idle; explicit waits |

## Next Steps

1. Review and approve these plans
2. Implement in order (Phase 1 → Phase 2 → Phase 3)
3. Run full test suite after each phase
4. Update documentation as needed
