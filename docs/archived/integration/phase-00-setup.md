# Integration Phase 0: API Client Refactoring

Prepare frontend to switch from mock to real APIs.

## Task 0.1: API Interface Definition

**Commits:**
- `refactor: define api interfaces in lib/api/types.ts`
- `refactor: update mock apis to implement interfaces`

**Example:**
```typescript
// lib/api/types.ts
export interface IAuthApi {
  login(email: string, password: string): Promise<User>;
  logout(): Promise<void>;
  getCurrentUser(): Promise<User | null>;
}

// lib/api/mock/auth.ts
export const mockAuthApi: IAuthApi = { ... };
```

## Task 0.2: Environment Configuration

**Commits:**
- `feat: add api mode environment variable`
- `feat: add api base url configuration`
- `feat: create api factory based on env`

**Example:**
```typescript
// lib/api/index.ts
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API === 'true';

export const authApi = USE_MOCK ? mockAuthApi : realAuthApi;
```

## Task 0.3: CORS and Proxy Setup

**Commits:**
- `chore: add vite proxy configuration for dev`
- `chore: update backend CORS settings`

---

## Total Integration Commits (Phase 0): 5
