package utils

import "go-gin-starter/models"

func HasPermission(role models.RoleEnum, permission string) bool {
	perms, exists := models.RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}
