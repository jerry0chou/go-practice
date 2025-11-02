package security

import (
	"errors"
	"fmt"
	"strings"
)

// Role represents a user role
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// Permission represents a system permission
type Permission struct {
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// User represents a user with roles and permissions
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
}

// RBACManager handles Role-Based Access Control
type RBACManager struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	users       map[string]*User
}

// NewRBACManager creates a new RBAC manager
func NewRBACManager() *RBACManager {
	return &RBACManager{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
		users:       make(map[string]*User),
	}
}

// AddPermission adds a new permission to the system
func (r *RBACManager) AddPermission(permission *Permission) {
	r.permissions[permission.Name] = permission
}

// AddRole adds a new role with permissions
func (r *RBACManager) AddRole(role *Role) error {
	// Validate that all permissions exist
	for _, permName := range role.Permissions {
		if _, exists := r.permissions[permName]; !exists {
			return fmt.Errorf("permission %s does not exist", permName)
		}
	}
	r.roles[role.Name] = role
	return nil
}

// AssignRoleToUser assigns a role to a user
func (r *RBACManager) AssignRoleToUser(userID, roleName string) error {
	if _, exists := r.roles[roleName]; !exists {
		return fmt.Errorf("role %s does not exist", roleName)
	}

	user, exists := r.users[userID]
	if !exists {
		return fmt.Errorf("user %s does not exist", userID)
	}

	// Check if user already has this role
	for _, existingRole := range user.Roles {
		if existingRole == roleName {
			return fmt.Errorf("user %s already has role %s", userID, roleName)
		}
	}

	user.Roles = append(user.Roles, roleName)
	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *RBACManager) RemoveRoleFromUser(userID, roleName string) error {
	user, exists := r.users[userID]
	if !exists {
		return fmt.Errorf("user %s does not exist", userID)
	}

	for i, role := range user.Roles {
		if role == roleName {
			user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user %s does not have role %s", userID, roleName)
}

// AddUser adds a new user to the system
func (r *RBACManager) AddUser(user *User) {
	r.users[user.ID] = user
}

// HasPermission checks if a user has a specific permission
func (r *RBACManager) HasPermission(userID, permissionName string) bool {
	user, exists := r.users[userID]
	if !exists {
		return false
	}

	// Check if user has any role that includes this permission
	for _, roleName := range user.Roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}

		for _, perm := range role.Permissions {
			if perm == permissionName {
				return true
			}
		}
	}

	return false
}

// HasRole checks if a user has a specific role
func (r *RBACManager) HasRole(userID, roleName string) bool {
	user, exists := r.users[userID]
	if !exists {
		return false
	}

	for _, role := range user.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}

// GetUserPermissions returns all permissions for a user
func (r *RBACManager) GetUserPermissions(userID string) ([]string, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, fmt.Errorf("user %s does not exist", userID)
	}

	var permissions []string
	permissionSet := make(map[string]bool)

	// Collect permissions from all user roles
	for _, roleName := range user.Roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}

		for _, perm := range role.Permissions {
			if !permissionSet[perm] {
				permissions = append(permissions, perm)
				permissionSet[perm] = true
			}
		}
	}

	return permissions, nil
}

// CheckResourceAccess checks if user can access a specific resource with an action
func (r *RBACManager) CheckResourceAccess(userID, resource, action string) bool {
	user, exists := r.users[userID]
	if !exists {
		return false
	}

	// Check if user has any role with permission for this resource and action
	for _, roleName := range user.Roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}

		for _, permName := range role.Permissions {
			permission, exists := r.permissions[permName]
			if !exists {
				continue
			}

			// Check if permission matches resource and action
			if permission.Resource == resource && permission.Action == action {
				return true
			}

			// Check for wildcard permissions
			if permission.Resource == "*" && permission.Action == action {
				return true
			}
			if permission.Resource == resource && permission.Action == "*" {
				return true
			}
			if permission.Resource == "*" && permission.Action == "*" {
				return true
			}
		}
	}

	return false
}

// GetUserRoles returns all roles for a user
func (r *RBACManager) GetUserRoles(userID string) ([]string, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, fmt.Errorf("user %s does not exist", userID)
	}

	return user.Roles, nil
}

// ValidatePermission validates if a permission string is properly formatted
func (r *RBACManager) ValidatePermission(permissionStr string) error {
	parts := strings.Split(permissionStr, ":")
	if len(parts) != 2 {
		return errors.New("permission must be in format 'resource:action'")
	}

	if parts[0] == "" || parts[1] == "" {
		return errors.New("permission resource and action cannot be empty")
	}

	return nil
}

// CreatePermissionFromString creates a permission from a string in format "resource:action"
func (r *RBACManager) CreatePermissionFromString(permissionStr, description string) (*Permission, error) {
	if err := r.ValidatePermission(permissionStr); err != nil {
		return nil, err
	}

	parts := strings.Split(permissionStr, ":")
	return &Permission{
		Name:        permissionStr,
		Resource:    parts[0],
		Action:      parts[1],
		Description: description,
	}, nil
}
