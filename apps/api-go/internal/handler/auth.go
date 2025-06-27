package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lucasfaria/rbac/api-go/internal/auth"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// CreateAccount handles user registration
// @Summary Create a new user account
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Param request body auth.CreateUserRequest true "User registration data"
// @Success 200 "User created successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /users [post]
func (h *AuthHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req auth.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	_, err := h.authService.CreateUser(req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AuthenticateWithPassword handles password-based authentication
// @Summary Authenticate with email and password
// @Description Authenticate user with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} auth.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /sessions/password [post]
func (h *AuthHandler) AuthenticateWithPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	authResponse, err := h.authService.AuthenticateWithPassword(req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

// AuthenticateWithGitHub handles GitHub OAuth authentication
// @Summary Authenticate with GitHub OAuth
// @Description Authenticate user with GitHub OAuth code, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.GitHubAuthRequest true "GitHub OAuth code"
// @Success 200 {object} auth.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /sessions/github [post]
func (h *AuthHandler) AuthenticateWithGitHub(w http.ResponseWriter, r *http.Request) {
	var req auth.GitHubAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Code == "" {
		http.Error(w, "GitHub OAuth code is required", http.StatusBadRequest)
		return
	}

	authResponse, err := h.authService.AuthenticateWithGitHub(req)
	if err != nil {
		http.Error(w, "GitHub authentication failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

// RegisterAuthRoutes registers authentication routes
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users", h.CreateAccount)
	r.Post("/sessions/password", h.AuthenticateWithPassword)
	r.Post("/sessions/github", h.AuthenticateWithGitHub)
}

type ErrorResponse struct {
	Error string `json:"error"`
}