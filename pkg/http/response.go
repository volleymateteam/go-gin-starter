package http

import (
	"fmt"
	"go-gin-starter/config"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	auth "go-gin-starter/pkg/auth"

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
		AvatarURL: fmt.Sprintf("https://%s/avatars/%s", config.AssetCloudFrontDomain, user.Avatar),
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
		"all_permissions":   auth.GetAllPermissions(user),
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
		"all_permissions":   auth.GetAllPermissions(user),
	}
}

// BuildTeamResponse builds TeamResponse DTO from Team model
func BuildTeamResponse(team *models.Team) dto.TeamResponse {
	logo := "defaults/default-team-logo.png"
	if team.Logo != "" {
		logo = team.Logo
	}

	return dto.TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		Country:   team.Country,
		SeasonID:  team.SeasonID,
		LogoURL:   fmt.Sprintf("https://%s/logos/teams/%s", config.AssetCloudFrontDomain, logo),
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
}
