package utils

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"mime/multipart"
	"path/filepath"
	"strings"
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

func ValidateImageFile(file *multipart.FileHeader) error {
	if file.Size > 2*1024*1024 {
		return errors.New(ErrLogoTooLarge)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return errors.New(ErrInvalidFileType)
	}
	return nil
}
