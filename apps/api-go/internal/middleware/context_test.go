package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test context helper functions without requiring database
func TestUserIDFromContext(t *testing.T) {
	t.Run("With valid user ID in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", "test-user-123")
		userID, ok := UserIDFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "test-user-123", userID)
	})

	t.Run("With no user ID in context", func(t *testing.T) {
		ctx := context.Background()
		userID, ok := UserIDFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, userID)
	})

	t.Run("With wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", 12345) // int instead of string
		userID, ok := UserIDFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, userID)
	})

	t.Run("With nil context value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", nil)
		userID, ok := UserIDFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, userID)
	})
}

func TestGetCurrentUserID(t *testing.T) {
	t.Run("With valid user ID in request context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), "userID", "test-user-456")
		req = req.WithContext(ctx)

		userID, err := GetCurrentUserID(req)
		assert.NoError(t, err)
		assert.Equal(t, "test-user-456", userID)
	})

	t.Run("With no user ID in request context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		
		userID, err := GetCurrentUserID(req)
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Equal(t, http.ErrNoCookie, err)
	})

	t.Run("With wrong type in request context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), "userID", 789) // int instead of string
		req = req.WithContext(ctx)

		userID, err := GetCurrentUserID(req)
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Equal(t, http.ErrNoCookie, err)
	})
}