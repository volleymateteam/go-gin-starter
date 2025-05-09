package utils

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
)

func ValidateAdminUpdateInput(input *dto.AdminUpdateUserInput) error {
	if input.Gender != "" && !models.IsValidGender(input.Gender) {
		return errors.New(ErrInvalidGender)
	}
	if input.Role != "" && !models.IsValidRole(input.Role) {
		return errors.New(ErrInvalidRole)
	}
	return nil
}
