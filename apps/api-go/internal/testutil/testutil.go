package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lucasfaria/rbac/api-go/internal/auth"
	"github.com/lucasfaria/rbac/api-go/internal/config"
	"github.com/lucasfaria/rbac/api-go/internal/database"
	db "github.com/lucasfaria/rbac/api-go/internal/database/sqlc"
	"github.com/stretchr/testify/require"
)

const (
	TestJWTSecret         = "test-jwt-secret-for-testing-only"
	TestGitHubClientID    = "test-github-client-id"
	TestGitHubSecret      = "test-github-client-secret"
	TestGitHubRedirectURI = "http://localhost:3000/test/callback"
)

// TestDB holds database connection and queries for testing
type TestDB struct {
	DB      *sql.DB
	Queries *db.Queries
}

// TestServices holds all services for testing
type TestServices struct {
	DB   *TestDB
	Auth *auth.AuthService
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	// Try to load config, fallback to environment variables
	cfg, err := config.LoadConfig(".")
	if err != nil {
		// Fallback to environment variables for CI/CD
		dbSource := os.Getenv("TEST_DB_SOURCE")
		if dbSource == "" {
			t.Skip("No test database configured. Set TEST_DB_SOURCE environment variable or create app.env file")
		}
		cfg.DB_SOURCE = dbSource
	}

	// Use the test database URL directly if TEST_DB_SOURCE is set
	testDBSource := cfg.DB_SOURCE
	if os.Getenv("TEST_DB_SOURCE") != "" {
		testDBSource = os.Getenv("TEST_DB_SOURCE")
	}
	
	dbConn, queries, err := database.NewDB(testDBSource)
	require.NoError(t, err, "Failed to connect to test database")

	testDB := &TestDB{
		DB:      dbConn,
		Queries: queries,
	}

	// Clean database before test
	CleanDatabase(t, testDB)

	return testDB
}

// SetupTestServices creates all services for testing
func SetupTestServices(t *testing.T) *TestServices {
	testDB := SetupTestDB(t)
	
	authService := auth.NewAuthService(
		testDB.DB,
		TestJWTSecret,
		TestGitHubClientID,
		TestGitHubSecret,
		TestGitHubRedirectURI,
	)

	return &TestServices{
		DB:   testDB,
		Auth: authService,
	}
}

// CleanDatabase removes all test data from the database
func CleanDatabase(t *testing.T, testDB *TestDB) {
	// Delete in reverse order due to foreign key constraints
	tables := []string{
		"accounts",
		"tokens",
		"members",
		"invites",
		"projects",
		"organizations",
		"users",
		"casbin_rule", // Casbin policy table
	}

	for _, table := range tables {
		_, err := testDB.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			// Some tables might not exist, that's okay
			log.Printf("Warning: Could not clean table %s: %v", table, err)
		}
	}
}

// TeardownTestDB closes the test database connection
func (tdb *TestDB) Teardown() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, testDB *TestDB, email, name, password string) *auth.User {
	authService := auth.NewAuthService(
		testDB.DB,
		TestJWTSecret,
		TestGitHubClientID,
		TestGitHubSecret,
		TestGitHubRedirectURI,
	)

	user, err := authService.CreateUser(auth.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	return user
}

// CreateTestOrganization creates a test organization in the database
func CreateTestOrganization(t *testing.T, testDB *TestDB, name, slug string, ownerID string, domain *string, shouldAttachUsersByDomain bool) {
	var domainValue sql.NullString
	if domain != nil {
		domainValue = sql.NullString{String: *domain, Valid: true}
	}

	_, err := testDB.DB.Exec(`
		INSERT INTO organizations (id, name, slug, domain, should_attach_users_by_domain, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, fmt.Sprintf("org-%s", slug), name, slug, domainValue, shouldAttachUsersByDomain, ownerID)
	require.NoError(t, err)
}

// CreateTestMembership creates a test membership in the database
func CreateTestMembership(t *testing.T, testDB *TestDB, userID, orgID, role string) {
	_, err := testDB.DB.Exec(`
		INSERT INTO members (id, role, organization_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, fmt.Sprintf("member-%s-%s", userID[:8], orgID[:8]), role, orgID, userID)
	require.NoError(t, err)
}

// ValidateJWTToken validates a JWT token and returns the user ID
func ValidateJWTToken(t *testing.T, token string) string {
	authService := auth.NewAuthService(
		nil, // DB not needed for token validation
		TestJWTSecret,
		TestGitHubClientID,
		TestGitHubSecret,
		TestGitHubRedirectURI,
	)

	userID, err := authService.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, userID)
	return userID
}