package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lucasfaria/rbac/api-go/internal/auth"
)

type AuthMiddleware struct {
	authService *auth.AuthService
}

func NewAuthMiddleware(authService *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// UserIDFromContext retrieves the user ID from the request context
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// RequireAuth is a middleware that requires authentication
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header must be Bearer token", http.StatusUnauthorized)
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Bearer token is required", http.StatusUnauthorized)
			return
		}

		// Validate token
		userID, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCurrentUserID is a helper function to get the current user ID from request context
func GetCurrentUserID(r *http.Request) (string, error) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		return "", http.ErrNoCookie // Using a standard error for "not authenticated"
	}
	return userID, nil
}

// GetUserMembership is a helper function to get user's membership in an organization
func (m *AuthMiddleware) GetUserMembership(userID, orgSlug string) (*auth.Membership, error) {
	return m.authService.GetUserMembership(userID, orgSlug)
}