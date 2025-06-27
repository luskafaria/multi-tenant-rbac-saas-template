package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJWTTokenGeneration tests JWT token generation without requiring database
func TestJWTTokenGeneration(t *testing.T) {
	// Create auth service with mock configuration (no DB needed for JWT operations)
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	userID := "test-user-123"

	// Generate token
	token, err := authService.GenerateToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	validatedUserID, err := authService.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestJWTTokenValidation(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	t.Run("Valid token", func(t *testing.T) {
		userID := "user-456"
		token, err := authService.GenerateToken(userID)
		require.NoError(t, err)

		validatedUserID, err := authService.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, validatedUserID)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		_, err := authService.ValidateToken("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing token")
	})

	t.Run("Empty token", func(t *testing.T) {
		_, err := authService.ValidateToken("")
		assert.Error(t, err)
	})

	t.Run("Token with wrong secret", func(t *testing.T) {
		// Create token with different secret
		wrongSecretService := &AuthService{
			jwtSecret: "wrong-secret",
		}
		token, err := wrongSecretService.GenerateToken("test-user")
		require.NoError(t, err)

		// Try to validate with correct secret service
		_, err = authService.ValidateToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing token")
	})
}

func TestJWTTokenExpiration(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	// Create an expired token manually for testing
	claims := Claims{
		UserID: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // Issued 2 hours ago
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(authService.jwtSecret))
	require.NoError(t, err)

	// Try to validate expired token
	_, err = authService.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token")
}

func TestJWTTokenStructure(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	userID := "test-user-789"
	token, err := authService.GenerateToken(userID)
	require.NoError(t, err)

	// Parse token to verify structure
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(authService.jwtSecret), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Check claims
	claims, ok := parsedToken.Claims.(*Claims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.NotZero(t, claims.IssuedAt)
	assert.NotZero(t, claims.ExpiresAt)
	
	// Verify expiration is ~7 days from now
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time
	timeDiff := actualExpiry.Sub(expectedExpiry).Abs()
	assert.Less(t, timeDiff, 1*time.Minute, "Token expiration should be approximately 7 days from now")
}