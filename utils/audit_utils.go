package utils

import (
	"go-gin-starter/dto"
	"go-gin-starter/models"
)

// BuildUserUpdateMetadata builds the audit log metadata for user update actions
func BuildUserUpdateMetadata(originalUser *models.User, input *dto.AdminUpdateUserInput) models.JSONBMap {
	updatedFields := []string{}
	metadata := models.JSONBMap{}

	if input.Username != "" && input.Username != originalUser.Username {
		updatedFields = append(updatedFields, "username")
		metadata["old_username"] = originalUser.Username
		metadata["new_username"] = input.Username
	}
	if input.Email != "" && input.Email != originalUser.Email {
		updatedFields = append(updatedFields, "email")
		metadata["old_email"] = originalUser.Email
		metadata["new_email"] = input.Email
	}
	if input.Gender != "" && input.Gender != originalUser.Gender {
		updatedFields = append(updatedFields, "gender")
		metadata["old_gender"] = originalUser.Gender
		metadata["new_gender"] = input.Gender
	}
	if input.Role != "" && input.Role != originalUser.Role {
		updatedFields = append(updatedFields, "role")
		metadata["old_role"] = originalUser.Role
		metadata["new_role"] = input.Role
	}

	metadata["updated_fields"] = updatedFields
	metadata["username"] = originalUser.Username
	metadata["email"] = originalUser.Email

	return metadata
}
