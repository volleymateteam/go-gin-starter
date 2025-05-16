package services

import (
	"errors"
	"fmt"
	"go-gin-starter/config"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	"go-gin-starter/pkg/storage"
	"go-gin-starter/repositories"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	CreateUser(username, email, password, gender string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUserByID(id uuid.UUID) error
	GetUsersWithPagination(page, limit int) ([]models.User, int64, error)
	AdminUpdateUser(id uuid.UUID, input *dto.AdminUpdateUserInput) (*models.User, error)
	UpdateUserPermissions(userID uuid.UUID, permissions []string) error
	GetUserProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateUserProfile(userID uuid.UUID, input dto.UpdateUserInput) error
	ChangeUserPassword(userID uuid.UUID, oldPassword, newPassword string) error
	DeleteUserProfile(userID uuid.UUID) error
	UploadUserAvatar(userID uuid.UUID, fileHeader *multipart.FileHeader) (string, error)
}

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser registers a new user
func (s *UserServiceImpl) CreateUser(username, email, password, gender string) (*models.User, error) {
	hashedPassword, err := authPkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Set default avatar based on gender
	var avatar string
	switch gender {
	case "male":
		avatar = "defaults/default-male.png"
	case "female":
		avatar = "defaults/default-female.png"
	default:
		avatar = "defaults/default-male.png"
	}

	var genderEnum models.GenderEnum
	switch gender {
	case "male":
		genderEnum = models.GenderMale
	case "female":
		genderEnum = models.GenderFemale
	default:
		genderEnum = models.GenderOther
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Gender:   genderEnum,
		Avatar:   avatar,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID fetches a user based on their UUID
func (s *UserServiceImpl) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// UpdateUser updates the user's information
func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUserByID soft-deletes a user based on UUID
func (s *UserServiceImpl) DeleteUserByID(id uuid.UUID) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(user)
}

// GetUsersWithPagination retrieves a paginated list of users
func (s *UserServiceImpl) GetUsersWithPagination(page int, limit int) ([]models.User, int64, error) {
	offset := (page - 1) * limit
	return s.userRepo.GetWithPagination(limit, offset)
}

// AdminUpdateUser updates user fields selectively
func (s *UserServiceImpl) AdminUpdateUser(id uuid.UUID, input *dto.AdminUpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserPermissions updates user permissions
func (s *UserServiceImpl) UpdateUserPermissions(userID uuid.UUID, permissions []string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.ExtraPermissions = models.StringArray(permissions) // Convert []string to StringArray

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// GetUserProfile retrieves detailed user profile information
func (s *UserServiceImpl) GetUserProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var avatarURL string
	if strings.HasPrefix(user.Avatar, "https://") {
		// Already a complete URL
		avatarURL = user.Avatar
	} else {
		// Default avatar or partial path, add CloudFront domain
		avatarURL = fmt.Sprintf("https://%s/avatars/%s", config.AssetCloudFrontDomain, user.Avatar)
	}

	// Get all permissions (combines role permissions with extra permissions)
	allPermissions := authPkg.GetAllPermissions(user)

	// Convert StringArray to []string for the DTO
	extraPermissions := []string{} // Initialize as empty array instead of nil/null
	if user.ExtraPermissions != nil {
		extraPermissions = []string(user.ExtraPermissions)
	}

	return &dto.UserResponse{
		ID:               user.ID,
		Username:         user.Username,
		Email:            user.Email,
		Gender:           string(user.Gender),
		AvatarURL:        avatarURL,
		Role:             string(user.Role),
		Permissions:      allPermissions,
		ExtraPermissions: extraPermissions,
		CreatedAt:        user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateUserProfile updates basic user information
func (s *UserServiceImpl) UpdateUserProfile(userID uuid.UUID, input dto.UpdateUserInput) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	return s.userRepo.Update(user)
}

// ChangeUserPassword updates the user's password
func (s *UserServiceImpl) ChangeUserPassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !authPkg.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New(constants.ErrPasswordMismatch)
	}

	hashed, err := authPkg.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	return s.userRepo.Update(user)
}

// DeleteUserProfile permanently deletes a user account
func (s *UserServiceImpl) DeleteUserProfile(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(user)
}

// UploadUserAvatar handles user avatar upload
func (s *UserServiceImpl) UploadUserAvatar(userID uuid.UUID, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > 2*1024*1024 {
		return "", errors.New(constants.ErrAvatarTooLarge)
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New(constants.ErrInvalidFileType)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", errors.New(constants.ErrUploadFailed)
	}
	defer src.Close()

	newFileName := uuid.New().String() + ext
	objectKey := fmt.Sprintf("avatars/%s", newFileName)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			contentType = "application/octet-stream"
		}
	}

	url, err := storage.UploadFileToS3(src, objectKey, contentType)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	user.Avatar = url
	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return url, nil
}
