package utils

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
