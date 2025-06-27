package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"

	"github.com/lucasfaria/rbac/api-go/internal/config"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Organization struct {
	ID                        string
	Name                      string
	Slug                      string
	Domain                    *string
	ShouldAttachUsersByDomain bool
	AvatarURL                 string
	OwnerID                   string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type Project struct {
	ID             string
	Name           string
	Description    string
	Slug           string
	AvatarURL      string
	OrganizationID string
	OwnerID        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Member struct {
	ID             string
	Role           string
	OrganizationID string
	UserID         string
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DB_SOURCE)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close()

	if err := clearDatabase(db); err != nil {
		log.Fatalf("failed to clear database: %v", err)
	}

	if err := seedDatabase(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}

func clearDatabase(db *sql.DB) error {
	// Clear in order due to foreign key constraints
	tables := []string{"members", "projects", "invites", "accounts", "tokens", "organizations", "users"}
	
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}
	
	fmt.Println("Database cleared")
	return nil
}

func seedDatabase(db *sql.DB) error {
	now := time.Now()
	
	// Create password hash for '123456'
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), 6)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create users
	users := []User{
		{
			ID:           uuid.New().String(),
			Name:         "John Doe",
			Email:        "john@gmail.com",
			PasswordHash: string(passwordHash),
			AvatarURL:    "https://github.com/luskafaria.png",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Name:         "Jane Smith",
			Email:        "jane@example.com",
			PasswordHash: string(passwordHash),
			AvatarURL:    "https://avatars.githubusercontent.com/u/1?v=4",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New().String(),
			Name:         "Bob Johnson",
			Email:        "bob@example.com",
			PasswordHash: string(passwordHash),
			AvatarURL:    "https://avatars.githubusercontent.com/u/2?v=4",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	for _, user := range users {
		_, err := db.Exec(`
			INSERT INTO users (id, name, email, password_hash, avatar_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, user.ID, user.Name, user.Email, user.PasswordHash, user.AvatarURL, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
	}
	fmt.Printf("Created %d users\n", len(users))

	// Create organizations
	acmeDomain := "acme.com"
	organizations := []Organization{
		{
			ID:                        uuid.New().String(),
			Name:                      "Acme Inc (Admin)",
			Slug:                      "acme-admin",
			Domain:                    &acmeDomain,
			ShouldAttachUsersByDomain: true,
			AvatarURL:                 "https://avatars.githubusercontent.com/u/3?v=4",
			OwnerID:                   users[0].ID,
			CreatedAt:                 now,
			UpdatedAt:                 now,
		},
		{
			ID:                        uuid.New().String(),
			Name:                      "Acme Inc (Member)",
			Slug:                      "acme-member",
			Domain:                    nil,
			ShouldAttachUsersByDomain: false,
			AvatarURL:                 "https://avatars.githubusercontent.com/u/4?v=4",
			OwnerID:                   users[0].ID,
			CreatedAt:                 now,
			UpdatedAt:                 now,
		},
		{
			ID:                        uuid.New().String(),
			Name:                      "Acme Inc (Billing)",
			Slug:                      "acme-billing",
			Domain:                    nil,
			ShouldAttachUsersByDomain: false,
			AvatarURL:                 "https://avatars.githubusercontent.com/u/5?v=4",
			OwnerID:                   users[0].ID,
			CreatedAt:                 now,
			UpdatedAt:                 now,
		},
	}

	for _, org := range organizations {
		_, err := db.Exec(`
			INSERT INTO organizations (id, name, slug, domain, should_attach_users_by_domain, avatar_url, owner_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, org.ID, org.Name, org.Slug, org.Domain, org.ShouldAttachUsersByDomain, org.AvatarURL, org.OwnerID, org.CreatedAt, org.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create organization %s: %w", org.Name, err)
		}
	}
	fmt.Printf("Created %d organizations\n", len(organizations))

	// Create members for each organization
	memberConfigs := []struct {
		orgIndex int
		roles    []string
	}{
		{0, []string{"ADMIN", "MEMBER", "MEMBER"}},   // Acme Admin
		{1, []string{"MEMBER", "ADMIN", "MEMBER"}},   // Acme Member  
		{2, []string{"BILLING", "ADMIN", "MEMBER"}},  // Acme Billing
	}

	memberCount := 0
	for _, config := range memberConfigs {
		org := organizations[config.orgIndex]
		for i, role := range config.roles {
			member := Member{
				ID:             uuid.New().String(),
				Role:           role,
				OrganizationID: org.ID,
				UserID:         users[i].ID,
			}

			_, err := db.Exec(`
				INSERT INTO members (id, role, organization_id, user_id)
				VALUES ($1, $2, $3, $4)
			`, member.ID, member.Role, member.OrganizationID, member.UserID)
			if err != nil {
				return fmt.Errorf("failed to create member: %w", err)
			}
			memberCount++
		}
	}
	fmt.Printf("Created %d members\n", memberCount)

	// Create projects for each organization
	projectCount := 0
	for _, org := range organizations {
		for i := 0; i < 3; i++ {
			project := Project{
				ID:             uuid.New().String(),
				Name:           fmt.Sprintf("Project %d for %s", i+1, org.Name),
				Description:    fmt.Sprintf("This is project %d for %s organization", i+1, org.Name),
				Slug:           fmt.Sprintf("project-%d-%s", i+1, org.Slug),
				AvatarURL:      fmt.Sprintf("https://avatars.githubusercontent.com/u/%d?v=4", 10+projectCount),
				OrganizationID: org.ID,
				OwnerID:        users[i%len(users)].ID, // Rotate between users
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			_, err := db.Exec(`
				INSERT INTO projects (id, name, description, slug, avatar_url, organization_id, owner_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, project.ID, project.Name, project.Description, project.Slug, project.AvatarURL, project.OrganizationID, project.OwnerID, project.CreatedAt, project.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to create project %s: %w", project.Name, err)
			}
			projectCount++
		}
	}
	fmt.Printf("Created %d projects\n", projectCount)

	return nil
}