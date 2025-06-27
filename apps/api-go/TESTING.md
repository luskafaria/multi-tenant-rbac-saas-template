# Testing Guide

This document describes the testing strategy and how to run tests for the RBAC API Go application.

## Test Structure

The project includes comprehensive integration tests covering:

### 1. Authentication Handler Tests (`internal/handler/auth_test.go`)
- **User Registration Tests:**
  - Successful user creation
  - Duplicate email validation
  - Input validation (missing fields, short passwords)
  - Auto-join organization functionality

- **Password Authentication Tests:**
  - Successful login with valid credentials
  - Invalid credentials handling (wrong password, wrong email, non-existent user)
  - Input validation (missing fields, invalid JSON)

### 2. GitHub OAuth Tests (`internal/handler/github_auth_test.go`)
- Request validation (missing code, invalid JSON)
- Invalid OAuth code handling
- *Note: Full OAuth flow testing requires mocking GitHub endpoints*

### 3. Authentication Middleware Tests (`internal/middleware/auth_test.go`)
- **JWT Token Validation:**
  - Successful authentication with valid token
  - Missing authorization header
  - Invalid header format (missing Bearer, empty token)
  - Invalid/expired JWT tokens

- **Context Helper Functions:**
  - `UserIDFromContext()` functionality
  - `GetCurrentUserID()` functionality

### 4. Integration Tests (`internal/integration/integration_test.go`)
- **Full Authentication Flow:**
  - User registration → Login → Access protected route
  - End-to-end authentication workflow

- **Protected Route Access:**
  - Access without authentication
  - Access with invalid tokens
  - Public route access (health check)

- **Error Scenarios:**
  - Duplicate user registration
  - Invalid login credentials

## Test Utilities (`internal/testutil/testutil.go`)

Provides helper functions for testing:
- **Database Setup:** Test database connection and cleanup
- **Service Setup:** Initialize all services with test configuration
- **Test Data Creation:** Create test users, organizations, memberships
- **JWT Validation:** Validate generated tokens

## Running Tests

### Prerequisites

1. **Test Database:** Ensure you have a test database configured
   ```bash
   # Set up test database URL in test.env or environment variable
   export TEST_DB_SOURCE="postgresql://docker:docker@localhost:5433/next-saas_test?sslmode=disable"
   ```

2. **Database Schema:** Make sure your test database has the required schema

### Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run only unit tests (excluding integration tests)
make test-unit

# Run only integration tests
make test-integration

# Run tests with race detection
make test-race

# Run specific test suites
go test -v ./internal/handler/... -run TestAuthHandlerTestSuite
go test -v ./internal/middleware/... -run TestAuthMiddlewareTestSuite
go test -v ./internal/integration/... -run TestIntegrationTestSuite
```

### Test Configuration

- **JWT Secret:** `test-jwt-secret-for-testing-only`
- **GitHub OAuth:** Mock configuration for testing
- **Database:** Uses separate test database with `_test` suffix

## Test Environment Variables

```bash
# test.env or environment variables
TEST_DB_SOURCE="postgresql://docker:docker@localhost:5433/next-saas_test?sslmode=disable"
API_PORT=3335
JWT_SECRET="test-jwt-secret-for-testing-only"
GITHUB_OAUTH_CLIENT_ID="test-github-client-id"
GITHUB_OAUTH_CLIENT_SECRET="test-github-client-secret"
GITHUB_OAUTH_CLIENT_REDIRECT_URI="http://localhost:3000/test/callback"
```

## Test Coverage

The tests cover:

✅ **Authentication Endpoints:**
- `POST /v1/users` (User registration)
- `POST /v1/sessions/password` (Password authentication)
- `POST /v1/sessions/github` (GitHub OAuth - validation only)

✅ **Authentication Middleware:**
- JWT token validation
- Authorization header parsing
- Context management

✅ **Protected Routes:**
- `GET /v1/profile` (Requires authentication)

✅ **Public Routes:**
- `GET /v1/health` (Health check)

✅ **Error Handling:**
- Input validation
- Authentication failures
- Authorization failures

## Continuous Integration

For CI/CD pipelines, ensure:

1. Test database is available
2. Environment variables are set
3. Run tests before deployment:
   ```bash
   make check  # Runs format, vet, and tests
   ```

## Extending Tests

When adding new features:

1. **Add unit tests** for individual components
2. **Add integration tests** for end-to-end flows
3. **Update test utilities** if new test data patterns are needed
4. **Maintain test isolation** - each test should clean up after itself

## Known Limitations

1. **GitHub OAuth:** Full OAuth flow testing requires mocking external services
2. **Database State:** Tests assume a clean database state (handled by test utilities)
3. **Time-based Tests:** JWT expiration testing may need time mocking for deterministic results

## Troubleshooting

### Common Issues

1. **Database Connection Failed:**
   - Verify test database is running
   - Check TEST_DB_SOURCE environment variable

2. **Tests Fail with "table doesn't exist":**
   - Ensure test database has the required schema
   - Run database migrations on test database

3. **Random Test Failures:**
   - Check for test isolation issues
   - Verify database cleanup between tests

### Debug Mode

Run tests with verbose output:
```bash
go test -v -run TestSpecificTest ./internal/handler/
```

Add debug logging in test utilities for troubleshooting database state.