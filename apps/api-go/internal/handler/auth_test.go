package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lucasfaria/rbac/api-go/internal/auth"
	"github.com/lucasfaria/rbac/api-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	services *testutil.TestServices
	handler  *AuthHandler
	router   *chi.Mux
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.services = testutil.SetupTestServices(suite.T())
	suite.handler = NewAuthHandler(suite.services.Auth)
	
	// Create router for testing
	suite.router = chi.NewRouter()
	suite.handler.RegisterRoutes(suite.router)
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	if suite.services != nil && suite.services.DB != nil {
		suite.services.DB.Teardown()
	}
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

// Test user registration
func (suite *AuthHandlerTestSuite) TestCreateAccount_Success() {
	// Prepare request
	reqBody := auth.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	// Make request
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	
	var user auth.User
	err = json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(suite.T(), err)
	
	assert.NotEmpty(suite.T(), user.ID)
	assert.Equal(suite.T(), "john@example.com", user.Email)
	assert.Equal(suite.T(), "John Doe", *user.Name)
	assert.NotZero(suite.T(), user.CreatedAt)
	assert.NotZero(suite.T(), user.UpdatedAt)
}

func (suite *AuthHandlerTestSuite) TestCreateAccount_DuplicateEmail() {
	// Create initial user
	testutil.CreateTestUser(suite.T(), suite.services.DB, "john@example.com", "John Doe", "password123")

	// Try to create user with same email
	reqBody := auth.CreateUserRequest{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Same email
		Password: "password456",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert conflict response
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "already exists")
}

func (suite *AuthHandlerTestSuite) TestCreateAccount_InvalidRequest() {
	testCases := []struct {
		name        string
		request     interface{}
		expectedMsg string
	}{
		{
			name: "Missing email",
			request: auth.CreateUserRequest{
				Name:     "John Doe",
				Password: "password123",
			},
			expectedMsg: "Email, password, and name are required",
		},
		{
			name: "Missing password",
			request: auth.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expectedMsg: "Email, password, and name are required",
		},
		{
			name: "Missing name",
			request: auth.CreateUserRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedMsg: "Email, password, and name are required",
		},
		{
			name: "Short password",
			request: auth.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "123", // Too short
			},
			expectedMsg: "Password must be at least 6 characters",
		},
		{
			name:        "Invalid JSON",
			request:     `{"invalid": json}`,
			expectedMsg: "Invalid request body",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			var jsonBody []byte
			var err error

			if str, ok := tc.request.(string); ok {
				jsonBody = []byte(str)
			} else {
				jsonBody, err = json.Marshal(tc.request)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedMsg)
		})
	}
}

func (suite *AuthHandlerTestSuite) TestCreateAccount_AutoJoinOrganization() {
	// Create an organization with domain auto-join
	ownerUser := testutil.CreateTestUser(suite.T(), suite.services.DB, "owner@company.com", "Owner", "password123")
	domain := "company.com"
	testutil.CreateTestOrganization(suite.T(), suite.services.DB, "Company Org", "company", ownerUser.ID, &domain, true)

	// Create user with matching domain
	reqBody := auth.CreateUserRequest{
		Name:     "New Employee",
		Email:    "employee@company.com", // Matching domain
		Password: "password123",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert user created successfully
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Verify user was auto-joined to organization
	var user auth.User
	err = json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(suite.T(), err)

	// Check membership exists
	var membershipCount int
	err = suite.services.DB.DB.QueryRow(`
		SELECT COUNT(*) FROM members m
		JOIN organizations o ON m.organization_id = o.id
		WHERE m.user_id = $1 AND o.slug = 'company'
	`, user.ID).Scan(&membershipCount)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, membershipCount)
}

// Test password authentication
func (suite *AuthHandlerTestSuite) TestAuthenticateWithPassword_Success() {
	// Create test user
	user := testutil.CreateTestUser(suite.T(), suite.services.DB, "john@example.com", "John Doe", "password123")

	// Prepare login request
	reqBody := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	// Make request
	req := httptest.NewRequest("POST", "/sessions/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var authResponse auth.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(suite.T(), err)
	
	// Validate token
	assert.NotEmpty(suite.T(), authResponse.Token)
	userID := testutil.ValidateJWTToken(suite.T(), authResponse.Token)
	assert.Equal(suite.T(), user.ID, userID)
	
	// Validate user data
	assert.Equal(suite.T(), user.ID, authResponse.User.ID)
	assert.Equal(suite.T(), user.Email, authResponse.User.Email)
	assert.Equal(suite.T(), user.Name, authResponse.User.Name)
}

func (suite *AuthHandlerTestSuite) TestAuthenticateWithPassword_InvalidCredentials() {
	// Create test user
	testutil.CreateTestUser(suite.T(), suite.services.DB, "john@example.com", "John Doe", "password123")

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "Wrong password",
			email:    "john@example.com",
			password: "wrongpassword",
		},
		{
			name:     "Wrong email",
			email:    "wrong@example.com",
			password: "password123",
		},
		{
			name:     "Non-existent user",
			email:    "nonexistent@example.com",
			password: "password123",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			reqBody := auth.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}
			jsonBody, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/sessions/password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid credentials")
		})
	}
}

func (suite *AuthHandlerTestSuite) TestAuthenticateWithPassword_InvalidRequest() {
	testCases := []struct {
		name        string
		request     interface{}
		expectedMsg string
	}{
		{
			name: "Missing email",
			request: auth.LoginRequest{
				Password: "password123",
			},
			expectedMsg: "Email and password are required",
		},
		{
			name: "Missing password",
			request: auth.LoginRequest{
				Email: "john@example.com",
			},
			expectedMsg: "Email and password are required",
		},
		{
			name:        "Invalid JSON",
			request:     `{"invalid": json}`,
			expectedMsg: "Invalid request body",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			var jsonBody []byte
			var err error

			if str, ok := tc.request.(string); ok {
				jsonBody = []byte(str)
			} else {
				jsonBody, err = json.Marshal(tc.request)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/sessions/password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedMsg)
		})
	}
}