package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test password hashing functionality without requiring database
func TestPasswordHashing(t *testing.T) {
	password := "testPassword123"
	
	// Test bcrypt hashing (this is what our CreateUser function uses)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	
	// Test password verification
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err, "Password should match its hash")
	
	// Test wrong password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongPassword"))
	assert.Error(t, err, "Wrong password should not match")
}

func TestPasswordHashingConsistency(t *testing.T) {
	password := "consistencyTest"
	
	// Generate multiple hashes of the same password
	hash1, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	require.NoError(t, err)
	
	hash2, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	require.NoError(t, err)
	
	// Hashes should be different (bcrypt includes salt)
	assert.NotEqual(t, hash1, hash2, "Bcrypt should generate different hashes for same password due to salt")
	
	// But both should verify correctly
	assert.NoError(t, bcrypt.CompareHashAndPassword(hash1, []byte(password)))
	assert.NoError(t, bcrypt.CompareHashAndPassword(hash2, []byte(password)))
}

func TestPasswordStrengthCosts(t *testing.T) {
	password := "strengthTest"
	
	testCases := []struct {
		name string
		cost int
	}{
		{"Low cost (fast)", 4},
		{"Default cost (medium)", 6}, // This is what we use in production
		{"High cost (slow)", 8},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), tc.cost)
			require.NoError(t, err)
			
			// Verify the password works
			err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
			assert.NoError(t, err)
			
			// Verify wrong password fails
			err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongPassword"))
			assert.Error(t, err)
		})
	}
}

func TestEmptyPasswordHandling(t *testing.T) {
	// Test empty password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(""), 6)
	require.NoError(t, err)
	
	// Empty password should verify correctly
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(""))
	assert.NoError(t, err)
	
	// Non-empty password should not match empty hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("notEmpty"))
	assert.Error(t, err)
}