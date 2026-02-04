# Testing Guide - Forge Admin Framework

## Overview

This guide provides a comprehensive testing strategy for the Forge full-stack admin framework with **maximum confidence and minimum tests**. The approach follows the test pyramid principle and focuses on critical user flows.

---

## Quick Start

### Install Dependencies

```bash
# Frontend
cd admin-ui
npm install

# Backend (already has testify)
cd ../forge
go mod download
```

### Run Tests

```bash
# Frontend unit tests
cd admin-ui
npm run test

# Backend tests
cd ../forge/admin/handlers
go test -v

# E2E tests
cd ../../admin-ui
npm run e2e
```

---

## Test Structure

### Test Pyramid (Total: 31 tests)

```
        E2E (5)           ← Few, slow, brittle
       /       \
    API (10)              ← Contract validation
   /           \
Unit (16)                 ← Many, fast, stable
```

- **Unit Tests**: 12 backend + 4 frontend = 16
- **API/Handler Tests**: 10
- **DB Integration Tests**: 4
- **E2E Tests**: 5

---

## Critical User Flows (Priority Order)

1. Admin Login → View List → Filter → View Details
2. Create Record with Validation
3. Edit Existing Record
4. Bulk Actions (Multi-select → Execute)
5. Relation Handling (Foreign Key Selection)
6. Pagination + Ordering
7. Permission-Based UI
8. Multi-Field Search

---

## Test Files

### Backend (Go)

```
forge/
├── admin/
│   ├── handlers/
│   │   └── list_test.go          ← API handler tests
│   └── core/
│       ├── admin_test.go          ← Unit tests (TODO)
│       └── permissions_test.go    ← Permission tests (TODO)
└── orm/
    └── queryset_test.go           ← DB integration tests (existing)
```

### Frontend (TypeScript)

```
admin-ui/
├── src/
│   ├── components/
│   │   ├── ModelList.test.tsx     ← Component tests
│   │   └── ModelForm.test.tsx     ← Form tests (TODO)
│   ├── api/
│   │   └── hooks/
│   │       └── adminHooks.test.ts ← Hook tests (TODO)
│   └── test/
│       └── setup.ts               ← Test setup
├── e2e/
│   └── admin-crud-flow.spec.ts    ← E2E tests
├── vitest.config.ts
└── playwright.config.ts
```

---

## Writing Tests

### Backend API Handler Test Pattern

```go
func TestHandler_Success(t *testing.T) {
    // 1. Setup
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    seedData(t, db)
    
    // 2. Create request
    req := httptest.NewRequest(http.MethodGet, "/api/admin/product/", nil)
    req = req.WithContext(context.WithValue(req.Context(), "user", testUser))
    
    // 3. Execute
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    
    // 4. Assert
    assert.Equal(t, http.StatusOK, w.Code)
    var response Response
    json.NewDecoder(w.Body).Decode(&response)
    assert.Equal(t, expectedValue, response.Field)
}
```

### Frontend Component Test Pattern

```typescript
it('renders data successfully', async () => {
  // 1. Mock dependencies
  vi.mocked(useModelList).mockReturnValue({
    data: mockData,
    isLoading: false,
    isError: false,
  });
  
  // 2. Render
  render(<ModelList model="product" />, { wrapper });
  
  // 3. Assert
  await waitFor(() => {
    expect(screen.getByRole('table')).toBeInTheDocument();
  });
  expect(screen.getByText('Product 1')).toBeInTheDocument();
});
```

### E2E Test Pattern

```typescript
test('complete CRUD flow', async ({ page }) => {
  // 1. Setup
  await page.request.post('/api/test/reset-db');
  
  // 2. Navigate
  await page.goto('/admin/product/');
  
  // 3. Interact
  await page.click('[data-testid="create-button"]');
  await page.fill('[name="name"]', 'Test Product');
  await page.click('[data-testid="submit-button"]');
  
  // 4. Assert
  await expect(page.locator('table')).toContainText('Test Product');
});
```

---

## Flakiness Prevention (8 Rules)

1. **Use `data-testid` selectors** - Never rely on text/CSS
2. **Seed deterministic data** - Fixed IDs, timestamps, sequential names
3. **Isolate tests** - Fresh DB per test (transactions or TRUNCATE)
4. **Mock time** - `vi.useFakeTimers()` or fixed `time.Now()`
5. **Wait for network idle** - `page.waitForLoadState('networkidle')`
6. **Disable animations** - `prefers-reduced-motion` in config
7. **Retry idempotent ops** - Use `waitFor` with timeout
8. **Avoid arbitrary waits** - Use conditions, not `sleep(1000)`

---

## What NOT to Test (6 Rules)

1. **Third-party internals** - Don't test TanStack Query caching
2. **CSS styling** - Avoid "button is blue" assertions
3. **Exact error wording** - Check presence, not text
4. **Every permutation** - Test one model deeply, assume others work
5. **Browser compatibility** - Chromium only is sufficient
6. **Responsive breakpoints** - One viewport (1280x720)

---

## Coverage Goals

- **Backend**: 70% line coverage
- **Frontend**: 60% line coverage
- **E2E**: Top 5-8 critical flows

**Remember**: Coverage is a metric, not a goal. Focus on critical paths.

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - run: go test -v -cover ./...
  
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npm run test:coverage
  
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npx playwright install --with-deps
      - run: npm run e2e
```

---

## Troubleshooting

### Tests are flaky
- Add `waitFor` around async assertions
- Use `data-testid` selectors
- Ensure DB cleanup between tests

### Database connection errors
- Check PostgreSQL is running
- Verify connection string
- Ensure test user has CREATE DATABASE permission

### Playwright browser errors
```bash
npx playwright install
```

### Module not found errors
```bash
cd admin-ui
npm install
```

---

## Next Steps

1. **Complete handler implementation** in `list_test.go`
2. **Add `data-testid` attributes** to frontend components
3. **Create test DB setup script** (`scripts/setup-test-db.sh`)
4. **Implement remaining priority tests** (see test_plan.md)
5. **Add CI/CD pipeline** (GitHub Actions)

---

## Resources

- [Full Test Plan](./test_plan.md) - Detailed test specifications
- [Implementation Summary](./test_implementation_summary.md) - Setup guide
- [Vitest Docs](https://vitest.dev/)
- [Playwright Docs](https://playwright.dev/)
- [React Testing Library](https://testing-library.com/react)
