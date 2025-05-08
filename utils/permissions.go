package utils

import (
	"go-gin-starter/models"
)

func HasPermission(role models.RoleEnum, permission string) bool {
	perms, exists := models.RolePermissions[role]
	if !exists {
		return false
	}

	// super admin shortcuts
	if contains(perms, "all") {
		return true
	}

	return contains(perms, permission)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
