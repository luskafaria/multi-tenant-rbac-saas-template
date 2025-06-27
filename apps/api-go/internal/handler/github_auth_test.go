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

type GitHubAuthTestSuite struct {
	suite.Suite
	services *testutil.TestServices
	handler  *AuthHandler
	router   *chi.Mux
}

func (suite *GitHubAuthTestSuite) SetupTest() {
	suite.services = testutil.SetupTestServices(suite.T())
	suite.handler = NewAuthHandler(suite.services.Auth)
	
	// Create router for testing
	suite.router = chi.NewRouter()
	suite.handler.RegisterRoutes(suite.router)
}

func (suite *GitHubAuthTestSuite) TearDownTest() {
	if suite.services != nil && suite.services.DB != nil {
		suite.services.DB.Teardown()
	}
}

func TestGitHubAuthTestSuite(t *testing.T) {
	suite.Run(t, new(GitHubAuthTestSuite))
}

// Test GitHub OAuth authentication validation
func (suite *GitHubAuthTestSuite) TestAuthenticateWithGitHub_InvalidRequest() {
	testCases := []struct {
		name        string
		request     interface{}
		expectedMsg string
	}{
		{
			name: "Missing code",
			request: auth.GitHubAuthRequest{
				Code: "",
			},
			expectedMsg: "GitHub OAuth code is required",
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

			req := httptest.NewRequest("POST", "/sessions/github", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedMsg)
		})
	}
}

func (suite *GitHubAuthTestSuite) TestAuthenticateWithGitHub_InvalidCode() {
	// Test with invalid GitHub OAuth code (will fail during token exchange)
	reqBody := auth.GitHubAuthRequest{
		Code: "invalid-github-code",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/sessions/github", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Should fail during OAuth token exchange
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "GitHub authentication failed")
}

// Note: Full GitHub OAuth flow testing would require mocking GitHub's OAuth endpoints
// For now, we test the endpoint validation and error handling
// In a real scenario, you would mock the HTTP client or use test doubles for the OAuth flow