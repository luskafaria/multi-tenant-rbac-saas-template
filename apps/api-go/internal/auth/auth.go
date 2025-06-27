package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthService struct {
	db           *sql.DB
	jwtSecret    string
	githubConfig *oauth2.Config
}

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        *string   `json:"name"`
	PasswordHash *string  `json:"-"`
	AvatarURL   *string   `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GitHubAuthRequest struct {
	Code string `json:"code"`
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, jwtSecret string, githubClientID, githubClientSecret, githubRedirectURI string) *AuthService {
	githubConfig := &oauth2.Config{
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		RedirectURL:  githubRedirectURI,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	return &AuthService{
		db:           db,
		jwtSecret:    jwtSecret,
		githubConfig: githubConfig,
	}
}

// CreateUser creates a new user account
func (a *AuthService) CreateUser(req CreateUserRequest) (*User, error) {
	// Check if user already exists
	var existingID string
	err := a.db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 6)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Generate UUID
	userID := uuid.New().String()

	// Insert user
	now := time.Now()
	_, err = a.db.Exec(`
		INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, req.Name, req.Email, string(hashedPassword), now, now)
	
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Check for auto-join organization
	if err := a.handleAutoJoinOrganization(userID, req.Email); err != nil {
		// Log error but don't fail user creation
		fmt.Printf("Error auto-joining organization: %v\n", err)
	}

	// Return created user
	user := &User{
		ID:        userID,
		Name:      &req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

// handleAutoJoinOrganization checks if user should auto-join an organization based on email domain
func (a *AuthService) handleAutoJoinOrganization(userID, email string) error {
	// Extract domain from email
	domain := ""
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			domain = email[i+1:]
			break
		}
	}

	if domain == "" {
		return nil // No domain found
	}

	// Find organization with matching domain that allows auto-join
	var orgID string
	err := a.db.QueryRow(`
		SELECT id FROM organizations 
		WHERE domain = $1 AND should_attach_users_by_domain = true
	`, domain).Scan(&orgID)
	
	if err == sql.ErrNoRows {
		return nil // No auto-join organization found
	}
	if err != nil {
		return fmt.Errorf("error finding auto-join organization: %w", err)
	}

	// Create membership
	memberID := uuid.New().String()
	_, err = a.db.Exec(`
		INSERT INTO members (id, role, organization_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, memberID, "MEMBER", orgID, userID)
	
	if err != nil {
		return fmt.Errorf("error creating membership: %w", err)
	}

	return nil
}

// AuthenticateWithPassword authenticates a user with email and password
func (a *AuthService) AuthenticateWithPassword(req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	var user User
	var passwordHash string
	err := a.db.QueryRow(`
		SELECT id, name, email, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := a.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GenerateToken generates a JWT token for a user
func (a *AuthService) GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

// ValidateToken validates a JWT token and returns the user ID
func (a *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token")
}

// GetUserByID gets a user by ID
func (a *AuthService) GetUserByID(userID string) (*User, error) {
	var user User
	err := a.db.QueryRow(`
		SELECT id, name, email, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

// GetUserMembership gets a user's membership in an organization
func (a *AuthService) GetUserMembership(userID, orgSlug string) (*Membership, error) {
	var membership Membership
	err := a.db.QueryRow(`
		SELECT m.id, m.role, m.organization_id, m.user_id
		FROM members m
		JOIN organizations o ON m.organization_id = o.id
		WHERE m.user_id = $1 AND o.slug = $2
	`, userID, orgSlug).Scan(&membership.ID, &membership.Role, &membership.OrganizationID, &membership.UserID)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("membership not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding membership: %w", err)
	}

	return &membership, nil
}

type Membership struct {
	ID             string `json:"id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
}

// AuthenticateWithGitHub authenticates a user with GitHub OAuth
func (a *AuthService) AuthenticateWithGitHub(req GitHubAuthRequest) (*AuthResponse, error) {
	// Exchange code for token
	token, err := a.githubConfig.Exchange(context.Background(), req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from GitHub
	client := a.githubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from GitHub: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GitHub response: %w", err)
	}

	var githubUser GitHubUser
	if err := json.Unmarshal(body, &githubUser); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	// Get user email if not provided in profile
	if githubUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, fmt.Errorf("failed to get user emails from GitHub: %w", err)
		}
		defer emailResp.Body.Close()

		emailBody, err := io.ReadAll(emailResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read GitHub emails response: %w", err)
		}

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err := json.Unmarshal(emailBody, &emails); err != nil {
			return nil, fmt.Errorf("failed to parse GitHub emails response: %w", err)
		}

		// Find primary email
		for _, email := range emails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}

	if githubUser.Email == "" {
		return nil, fmt.Errorf("no email found in GitHub profile")
	}

	// Find or create user
	user, err := a.findOrCreateGitHubUser(githubUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token
	jwtToken, err := a.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &AuthResponse{
		Token: jwtToken,
		User:  *user,
	}, nil
}

// findOrCreateGitHubUser finds an existing user or creates a new one for GitHub authentication
func (a *AuthService) findOrCreateGitHubUser(githubUser GitHubUser) (*User, error) {
	// First, check if account already exists
	var userID string
	err := a.db.QueryRow(`
		SELECT user_id FROM accounts WHERE provider = 'GITHUB' AND provider_account_id = $1
	`, fmt.Sprintf("%d", githubUser.ID)).Scan(&userID)

	if err == nil {
		// Account exists, get user
		return a.GetUserByID(userID)
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing account: %w", err)
	}

	// Check if user exists by email
	var existingUser User
	err = a.db.QueryRow(`
		SELECT id, name, email, avatar_url, created_at, updated_at
		FROM users WHERE email = $1
	`, githubUser.Email).Scan(&existingUser.ID, &existingUser.Name, &existingUser.Email, 
		&existingUser.AvatarURL, &existingUser.CreatedAt, &existingUser.UpdatedAt)

	if err == nil {
		// User exists, create account link
		accountID := uuid.New().String()
		_, err = a.db.Exec(`
			INSERT INTO accounts (id, provider, provider_account_id, user_id)
			VALUES ($1, $2, $3, $4)
		`, accountID, "GITHUB", fmt.Sprintf("%d", githubUser.ID), existingUser.ID)
		
		if err != nil {
			return nil, fmt.Errorf("error linking GitHub account: %w", err)
		}

		// Update avatar URL if not set
		if existingUser.AvatarURL == nil && githubUser.AvatarURL != "" {
			_, err = a.db.Exec(`
				UPDATE users SET avatar_url = $1, updated_at = $2 WHERE id = $3
			`, githubUser.AvatarURL, time.Now(), existingUser.ID)
			if err != nil {
				// Log but don't fail
				fmt.Printf("Error updating avatar URL: %v\n", err)
			} else {
				existingUser.AvatarURL = &githubUser.AvatarURL
			}
		}

		return &existingUser, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	// Create new user
	userID = uuid.New().String()
	now := time.Now()
	
	name := githubUser.Name
	if name == "" {
		name = githubUser.Login
	}

	_, err = a.db.Exec(`
		INSERT INTO users (id, name, email, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, name, githubUser.Email, githubUser.AvatarURL, now, now)
	
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Create account link
	accountID := uuid.New().String()
	_, err = a.db.Exec(`
		INSERT INTO accounts (id, provider, provider_account_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, accountID, "GITHUB", fmt.Sprintf("%d", githubUser.ID), userID)
	
	if err != nil {
		return nil, fmt.Errorf("error creating GitHub account link: %w", err)
	}

	// Handle auto-join organization
	if err := a.handleAutoJoinOrganization(userID, githubUser.Email); err != nil {
		fmt.Printf("Error auto-joining organization: %v\n", err)
	}

	// Return created user
	user := &User{
		ID:        userID,
		Name:      &name,
		Email:     githubUser.Email,
		AvatarURL: &githubUser.AvatarURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}