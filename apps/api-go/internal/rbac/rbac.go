package rbac

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RBAC struct {
	enforcer *casbin.Enforcer
}

// NewRBAC creates a new RBAC instance with casbin enforcer
func NewRBAC(databaseURL string) (*RBAC, error) {
	// Initialize GORM with PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// Load model from file
	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Enable auto-save
	enforcer.EnableAutoSave(true)

	// Initialize policies
	rbac := &RBAC{enforcer: enforcer}
	if err := rbac.initializePolicies(); err != nil {
		return nil, fmt.Errorf("failed to initialize policies: %w", err)
	}

	return rbac, nil
}

// initializePolicies sets up the default RBAC policies based on the roles from the original system
func (r *RBAC) initializePolicies() error {
	// Clear existing policies
	r.enforcer.ClearPolicy()

	// Define role-based policies
	policies := [][]string{
		// ADMIN role permissions
		{"role:admin", "User", "get"},
		{"role:admin", "User", "create"},
		{"role:admin", "User", "update"},
		{"role:admin", "User", "delete"},
		{"role:admin", "Organization", "get"},
		{"role:admin", "Organization", "create"},
		{"role:admin", "Organization", "update"},
		{"role:admin", "Organization", "delete"},
		{"role:admin", "Organization", "manage"},
		{"role:admin", "Project", "get"},
		{"role:admin", "Project", "create"},
		{"role:admin", "Project", "update"},
		{"role:admin", "Project", "delete"},
		{"role:admin", "Project", "manage"},
		{"role:admin", "Invite", "get"},
		{"role:admin", "Invite", "create"},
		{"role:admin", "Invite", "update"},
		{"role:admin", "Invite", "delete"},
		{"role:admin", "Member", "get"},
		{"role:admin", "Member", "create"},
		{"role:admin", "Member", "update"},
		{"role:admin", "Member", "delete"},
		{"role:admin", "Billing", "get"},
		{"role:admin", "Billing", "update"},

		// MEMBER role permissions (limited access)
		{"role:member", "User", "get"},
		{"role:member", "Organization", "get"},
		{"role:member", "Project", "get"},
		{"role:member", "Project", "create"},
		{"role:member", "Member", "get"},

		// BILLING role permissions (only billing access)
		{"role:billing", "Billing", "get"},
		{"role:billing", "Billing", "update"},
	}

	// Add policies
	for _, policy := range policies {
		if added, err := r.enforcer.AddPolicy(policy); err != nil {
			log.Printf("Error adding policy %v: %v", policy, err)
		} else if !added {
			log.Printf("Policy already exists: %v", policy)
		}
	}

	// Save policies
	return r.enforcer.SavePolicy()
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (r *RBAC) CheckPermission(userID, orgID, resource, action string) bool {
	subject := fmt.Sprintf("user:%s:org:%s", userID, orgID)
	result, err := r.enforcer.Enforce(subject, resource, action)
	if err != nil {
		// Log error and deny access on error
		fmt.Printf("Error checking permission: %v\n", err)
		return false
	}
	return result
}

// AssignRole assigns a role to a user in an organization
func (r *RBAC) AssignRole(userID, orgID, role string) error {
	subject := fmt.Sprintf("user:%s:org:%s", userID, orgID)
	roleSubject := fmt.Sprintf("role:%s", role)

	added, err := r.enforcer.AddGroupingPolicy(subject, roleSubject)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}
	if !added {
		log.Printf("Role assignment already exists: user %s, org %s, role %s", userID, orgID, role)
	}

	return nil
}

// RemoveRole removes a role from a user in an organization
func (r *RBAC) RemoveRole(userID, orgID, role string) error {
	subject := fmt.Sprintf("user:%s:org:%s", userID, orgID)
	roleSubject := fmt.Sprintf("role:%s", role)

	removed, err := r.enforcer.RemoveGroupingPolicy(subject, roleSubject)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	if !removed {
		log.Printf("Role assignment did not exist: user %s, org %s, role %s", userID, orgID, role)
	}

	return nil
}

// GetUserRoles gets all roles for a user in an organization
func (r *RBAC) GetUserRoles(userID, orgID string) []string {
	subject := fmt.Sprintf("user:%s:org:%s", userID, orgID)
	roles, _ := r.enforcer.GetRolesForUser(subject)

	// Clean up role names (remove "role:" prefix)
	cleanRoles := make([]string, len(roles))
	for i, role := range roles {
		if len(role) > 5 && role[:5] == "role:" {
			cleanRoles[i] = role[5:]
		} else {
			cleanRoles[i] = role
		}
	}

	return cleanRoles
}

// CheckResourceOwnership checks if a user owns a specific resource
func (r *RBAC) CheckResourceOwnership(userID, resourceType, resourceID string) bool {
	// For ownership checks, we use a special permission format
	subject := fmt.Sprintf("user:%s", userID)
	resource := fmt.Sprintf("%s:%s", resourceType, resourceID)
	result, err := r.enforcer.Enforce(subject, resource, "owner")
	if err != nil {
		// Log error and deny access on error
		fmt.Printf("Error checking resource ownership: %v\n", err)
		return false
	}
	return result
}

// AssignResourceOwnership assigns ownership of a resource to a user
func (r *RBAC) AssignResourceOwnership(userID, resourceType, resourceID string) error {
	subject := fmt.Sprintf("user:%s", userID)
	resource := fmt.Sprintf("%s:%s", resourceType, resourceID)

	added, err := r.enforcer.AddPolicy(subject, resource, "owner")
	if err != nil {
		return fmt.Errorf("failed to assign ownership: %w", err)
	}
	if !added {
		log.Printf("Ownership already exists: user %s, resource %s:%s", userID, resourceType, resourceID)
	}

	return nil
}

// RemoveResourceOwnership removes ownership of a resource from a user
func (r *RBAC) RemoveResourceOwnership(userID, resourceType, resourceID string) error {
	subject := fmt.Sprintf("user:%s", userID)
	resource := fmt.Sprintf("%s:%s", resourceType, resourceID)

	removed, err := r.enforcer.RemovePolicy(subject, resource, "owner")
	if err != nil {
		return fmt.Errorf("failed to remove ownership: %w", err)
	}
	if !removed {
		log.Printf("Ownership did not exist: user %s, resource %s:%s", userID, resourceType, resourceID)
	}

	return nil
}
