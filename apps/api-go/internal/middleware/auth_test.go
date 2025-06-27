package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lucasfaria/rbac/api-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	services   *testutil.TestServices
	middleware *AuthMiddleware
	router     *chi.Mux
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.services = testutil.SetupTestServices(suite.T())
	suite.middleware = NewAuthMiddleware(suite.services.Auth)
	
	// Create router with middleware for testing
	suite.router = chi.NewRouter()
	
	// Add protected route for testing
	suite.router.Group(func(r chi.Router) {
		r.Use(suite.middleware.RequireAuth)
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "No user ID in context", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"user_id":"` + userID + `","message":"authenticated"}`))
		})
	})
}

func (suite *AuthMiddlewareTestSuite) TearDownTest() {
	if suite.services != nil && suite.services.DB != nil {
		suite.services.DB.Teardown()
	}
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

// Test successful authentication
func (suite *AuthMiddlewareTestSuite) TestRequireAuth_Success() {
	// Create test user and generate token
	user := testutil.CreateTestUser(suite.T(), suite.services.DB, "john@example.com", "John Doe", "password123")
	token, err := suite.services.Auth.GenerateToken(user.ID)
	require.NoError(suite.T(), err)

	// Make request with valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert successful response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), user.ID)
	assert.Contains(suite.T(), w.Body.String(), "authenticated")
}

// Test missing authorization header
func (suite *AuthMiddlewareTestSuite) TestRequireAuth_MissingHeader() {
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Authorization header required")
}

// Test invalid authorization header format
func (suite *AuthMiddlewareTestSuite) TestRequireAuth_InvalidHeaderFormat() {
	testCases := []struct {
		name   string
		header string
		msg    string
	}{
		{
			name:   "Missing Bearer prefix",
			header: "invalid-token",
			msg:    "Authorization header must be Bearer token",
		},
		{
			name:   "Empty Bearer token",
			header: "Bearer ",
			msg:    "Bearer token is required",
		},
		{
			name:   "Only Bearer",
			header: "Bearer",
			msg:    "Authorization header must be Bearer token",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), tc.msg)
		})
	}
}

// Test invalid JWT token
func (suite *AuthMiddlewareTestSuite) TestRequireAuth_InvalidToken() {
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Malformed token",
			token: "invalid.jwt.token",
		},
		{
			name:  "Expired token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjJ9.invalid", // Expired
		},
		{
			name:  "Random string",
			token: "not-a-jwt-token-at-all",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid or expired token")
		})
	}
}

// Test user ID context helper functions
func (suite *AuthMiddlewareTestSuite) TestUserIDFromContext() {
	// Test with user ID in context
	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	userID, ok := UserIDFromContext(ctx)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "test-user-id", userID)

	// Test with no user ID in context
	ctx = context.Background()
	userID, ok = UserIDFromContext(ctx)
	assert.False(suite.T(), ok)
	assert.Empty(suite.T(), userID)

	// Test with wrong type in context
	ctx = context.WithValue(context.Background(), "userID", 123)
	userID, ok = UserIDFromContext(ctx)
	assert.False(suite.T(), ok)
	assert.Empty(suite.T(), userID)
}

func (suite *AuthMiddlewareTestSuite) TestGetCurrentUserID() {
	// Test with user ID in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), "userID", "test-user-id")
	req = req.WithContext(ctx)

	userID, err := GetCurrentUserID(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-user-id", userID)

	// Test with no user ID in context
	req = httptest.NewRequest("GET", "/test", nil)
	userID, err = GetCurrentUserID(req)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), userID)
}