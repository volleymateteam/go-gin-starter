package validation

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func ValidateAdminUpdateInput(input *dto.AdminUpdateUserInput) error {
	if input.Gender != "" && !models.IsValidGender(input.Gender) {
		return errors.New(constants.ErrInvalidGender)
	}
	if input.Role != "" && !models.IsValidRole(input.Role) {
		return errors.New(constants.ErrInvalidRole)
	}
	return nil
}

func ValidateImageFile(file *multipart.FileHeader) error {
	if file.Size > 2*1024*1024 {
		return errors.New(constants.ErrLogoTooLarge)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return errors.New(constants.ErrInvalidFileType)
	}
	return nil
}
