package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lucasfaria/rbac/api-go/internal/auth"
	"github.com/stretchr/testify/assert"
)

// Test handler validation logic without database dependency
func TestCreateAccountValidation(t *testing.T) {
	// Create a mock auth service that will never be called due to validation failures
	var mockAuthService *auth.AuthService = nil
	handler := NewAuthHandler(mockAuthService)
	
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	testCases := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing email",
			requestBody: map[string]string{
				"name":     "John Doe",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email, password, and name are required",
		},
		{
			name: "Missing password",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email, password, and name are required",
		},
		{
			name: "Missing name",
			requestBody: map[string]string{
				"email":    "john@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email, password, and name are required",
		},
		{
			name: "Short password",
			requestBody: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "123", // Too short
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password must be at least 6 characters",
		},
		{
			name: "Empty fields",
			requestBody: map[string]string{
				"name":     "",
				"email":    "",
				"password": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email, password, and name are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestAuthenticateWithPasswordValidation(t *testing.T) {
	// Create a mock auth service that will never be called due to validation failures
	var mockAuthService *auth.AuthService = nil
	handler := NewAuthHandler(mockAuthService)
	
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	testCases := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing email",
			requestBody: map[string]string{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email and password are required",
		},
		{
			name: "Missing password",
			requestBody: map[string]string{
				"email": "john@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email and password are required",
		},
		{
			name: "Empty fields",
			requestBody: map[string]string{
				"email":    "",
				"password": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email and password are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/sessions/password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestGitHubAuthValidation(t *testing.T) {
	// Create a mock auth service that will never be called due to validation failures
	var mockAuthService *auth.AuthService = nil
	handler := NewAuthHandler(mockAuthService)
	
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	testCases := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing code",
			requestBody: map[string]string{
				"code": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "GitHub OAuth code is required",
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "GitHub OAuth code is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/sessions/github", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestInvalidJSONHandling(t *testing.T) {
	var mockAuthService *auth.AuthService = nil
	handler := NewAuthHandler(mockAuthService)
	
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	testCases := []struct {
		name     string
		endpoint string
		body     string
	}{
		{
			name:     "Invalid JSON for user creation",
			endpoint: "/users",
			body:     `{"name": "John", "email": invalid json}`,
		},
		{
			name:     "Invalid JSON for password auth",
			endpoint: "/sessions/password",
			body:     `{"email": "test@example.com", invalid}`,
		},
		{
			name:     "Invalid JSON for GitHub auth",
			endpoint: "/sessions/github",
			body:     `{"code": incomplete json`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.endpoint, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid request body")
		})
	}
}