package auth

import (
	"go-gin-starter/models"
)

func HasPermission(user *models.User, permission string) bool {
	// 1. Check role-based permissions
	rolePerms, exists := models.RolePermissions[user.Role]
	if exists {
		if contains(rolePerms, "all") || contains(rolePerms, permission) {
			return true
		}
	}

	// 2. Check extra per-user permissions
	if contains(user.ExtraPermissions, permission) {
		return true
	}

	// 3. Deny by default
	return false
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// GetAllPermissions returns all effective permissions (role + extra)
func GetAllPermissions(user *models.User) []string {
	// Start with role permissions
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	// Create a map for quick lookup and to avoid duplicates
	permMap := make(map[string]bool)

	// Add role permissions
	for _, p := range rolePerms {
		permMap[p] = true
	}

	// Add extra permissions
	for _, p := range user.ExtraPermissions {
		permMap[p] = true
	}

	// Convert back to slice
	result := make([]string, 0, len(permMap))
	for p := range permMap {
		result = append(result, p)
	}

	return result
}
