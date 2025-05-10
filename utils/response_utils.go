package utils

import (
	"go-gin-starter/dto"
	"go-gin-starter/models"

	"github.com/gin-gonic/gin"
)

// BuildAdminUserResponse constructs AdminUserResponse DTO from User model
func BuildAdminUserResponse(user *models.User) dto.AdminUserResponse {
	return dto.AdminUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Gender:    user.Gender,
		Role:      user.Role,
		AvatarURL: "/uploads/avatars/" + user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// BuildUserPermissionsResponse builds the API response for user permissions update
func BuildUserPermissionsResponse(user *models.User) gin.H {
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	return gin.H{
		"user_id":           user.ID,
		"username":          user.Username,
		"role":              user.Role,
		"role_permissions":  rolePerms,
		"extra_permissions": user.ExtraPermissions,
		"all_permissions":   GetAllPermissions(user),
	}
}

// BuildUserResetPermissionsResponse builds API response for resetting user permissions
func BuildUserResetPermissionsResponse(user *models.User, emptyPermissions []string) gin.H {
	rolePerms, exists := models.RolePermissions[user.Role]
	if !exists {
		rolePerms = []string{}
	}

	return gin.H{
		"user_id":           user.ID,
		"username":          user.Username,
		"role":              user.Role,
		"role_permissions":  rolePerms,
		"extra_permissions": []string{},
		"all_permissions":   GetAllPermissions(user),
	}
}

// BuildTeamResponse builds TeamResponse DTO from Team model
func BuildTeamResponse(team *models.Team) dto.TeamResponse {
	logoPath := "/uploads/logos/defaults/default-team-logo.png"
	if team.Logo != "" {
		logoPath = "/uploads/logos/" + team.Logo
	}

	return dto.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Country:   team.Country,
		SeasonID:  team.SeasonID,
		LogoURL:   logoPath,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
}
