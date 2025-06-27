package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucasfaria/rbac/api-go/internal/auth"
	"github.com/lucasfaria/rbac/api-go/internal/router"
	"github.com/lucasfaria/rbac/api-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	services *testutil.TestServices
	handler  http.Handler
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.services = testutil.SetupTestServices(suite.T())
	
	// Create full router with all dependencies
	suite.handler = router.NewRouter(
		suite.services.DB.Queries,
		suite.services.DB.DB,
		testutil.TestJWTSecret,
		testutil.TestGitHubClientID,
		testutil.TestGitHubSecret,
		testutil.TestGitHubRedirectURI,
	)
}

func (suite *IntegrationTestSuite) TearDownTest() {
	if suite.services != nil && suite.services.DB != nil {
		suite.services.DB.Teardown()
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Test full authentication flow
func (suite *IntegrationTestSuite) TestFullAuthenticationFlow() {
	// Step 1: Register a new user
	registerReq := auth.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, err := json.Marshal(registerReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	
	var user auth.User
	err = json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "john@example.com", user.Email)

	// Step 2: Login with the created user
	loginReq := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, err = json.Marshal(loginReq)
	require.NoError(suite.T(), err)

	req = httptest.NewRequest("POST", "/v1/sessions/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var authResponse auth.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), authResponse.Token)
	assert.Equal(suite.T(), user.ID, authResponse.User.ID)

	// Step 3: Access protected route with token
	req = httptest.NewRequest("GET", "/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	w = httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	// The profile endpoint should return user data
	assert.Contains(suite.T(), w.Body.String(), user.ID)
}

// Test accessing protected route without authentication
func (suite *IntegrationTestSuite) TestProtectedRouteWithoutAuth() {
	req := httptest.NewRequest("GET", "/v1/profile", nil)
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Authorization header required")
}

// Test accessing protected route with invalid token
func (suite *IntegrationTestSuite) TestProtectedRouteWithInvalidToken() {
	req := httptest.NewRequest("GET", "/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid or expired token")
}

// Test health check endpoint (public route)
func (suite *IntegrationTestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/v1/health", nil)
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "OK")
}

// Test user registration with duplicate email
func (suite *IntegrationTestSuite) TestDuplicateUserRegistration() {
	// Create first user
	registerReq := auth.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, err := json.Marshal(registerReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Try to create user with same email
	registerReq2 := auth.CreateUserRequest{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Same email
		Password: "password456",
	}
	jsonBody, err = json.Marshal(registerReq2)
	require.NoError(suite.T(), err)

	req = httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "already exists")
}

// Test invalid login credentials
func (suite *IntegrationTestSuite) TestInvalidLoginCredentials() {
	// Create user first
	user := testutil.CreateTestUser(suite.T(), suite.services.DB, "john@example.com", "John Doe", "password123")
	require.NotNil(suite.T(), user)

	// Try to login with wrong password
	loginReq := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	jsonBody, err := json.Marshal(loginReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/v1/sessions/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid credentials")
}