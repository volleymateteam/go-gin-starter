package services

import (
	"errors"
	"go-gin-starter/models"
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	"go-gin-starter/repositories"
	"time"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(username, email, password, gender string) (*models.User, error)
	Login(email, password string) (string, string, error)     // Returns access token, refresh token, error
	RefreshToken(refreshToken string) (string, string, error) // Returns new access token, new refresh token, error
	ForgotPassword(email string) (string, error)              // Returns reset token, error
	ResetPassword(token, newPassword string) error
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	authRepo repositories.AuthRepository
	userRepo repositories.UserRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repositories.AuthRepository, userRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

// Register handles user registration
func (s *AuthServiceImpl) Register(username, email, password, gender string) (*models.User, error) {
	// Check if password is strong
	if !authPkg.IsStrongPassword(password) {
		return nil, errors.New(constants.ErrStrongPassword)
	}

	// Check if both username and email already exist
	usernameExists := false
	emailExists := false

	// Check email
	existingUser, err := s.authRepo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		emailExists = true
	}

	// Check username
	existingUser, err = s.authRepo.GetUserByUsername(username)
	if err == nil && existingUser != nil {
		usernameExists = true
	}

	// Return appropriate error based on what exists
	if emailExists && usernameExists {
		return nil, errors.New(constants.ErrBothUserAndEmailExist)
	} else if emailExists {
		return nil, errors.New(constants.ErrEmailAlreadyExists)
	} else if usernameExists {
		return nil, errors.New(constants.ErrUserAlreadyExists)
	}

	// Only allow "male" or "female" genders for avatar and genderEnum
	if gender != "male" && gender != "female" {
		return nil, errors.New(constants.ErrInvalidGender)
	}

	// Set default avatar based on gender
	var avatar string
	if gender == "male" {
		avatar = "defaults/default-male.png"
	} else {
		avatar = "defaults/default-female.png"
	}

	var genderEnum models.GenderEnum
	if gender == "male" {
		genderEnum = models.GenderMale
	} else {
		genderEnum = models.GenderFemale
	}

	// Hash the password
	hashedPassword, err := authPkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
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

// Login handles user authentication
func (s *AuthServiceImpl) Login(email, password string) (string, string, error) {
	// Find the user by email
	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New(constants.ErrInvalidCredentials)
	}

	// Check password
	if !authPkg.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New(constants.ErrInvalidCredentials)
	}

	// Generate access token
	accessToken, err := authPkg.GenerateJWT(user.ID)
	if err != nil {
		return "", "", errors.New(constants.ErrTokenGenerationFailed)
	}

	// Generate refresh token
	refreshToken, err := authPkg.GenerateSecureToken(32)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	// Set refresh token expiry
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Update user with refresh token
	if err := s.authRepo.UpdateRefreshToken(user.ID, refreshToken, refreshExpiry); err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

// RefreshToken handles refreshing the JWT using a refresh token
func (s *AuthServiceImpl) RefreshToken(refreshToken string) (string, string, error) {
	// Find the user by refresh token
	user, err := s.authRepo.GetUserByRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New(constants.ErrInvalidToken)
	}

	// Check if refresh token is expired
	if user.RefreshTokenExpiry.Before(time.Now()) {
		return "", "", errors.New(constants.ErrInvalidToken)
	}

	// Generate new access token
	accessToken, err := authPkg.GenerateJWT(user.ID)
	if err != nil {
		return "", "", errors.New(constants.ErrTokenGenerationFailed)
	}

	// Generate new refresh token
	newRefreshToken, err := authPkg.GenerateSecureToken(32)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	// Set new refresh token expiry
	newExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Update user with new refresh token
	if err := s.authRepo.UpdateRefreshToken(user.ID, newRefreshToken, newExpiry); err != nil {
		return "", "", errors.New("failed to update refresh token")
	}

	return accessToken, newRefreshToken, nil
}

// ForgotPassword handles password reset request
func (s *AuthServiceImpl) ForgotPassword(email string) (string, error) {
	// Check if user exists
	_, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New(constants.ErrUserNotFound)
	}

	// Generate reset token
	resetToken, err := authPkg.GenerateSecureToken(16)
	if err != nil {
		return "", errors.New(constants.ErrResetTokenFailed)
	}

	// Set token expiry (12 hours)
	expiry := time.Now().Add(12 * time.Hour)

	// Update user with reset token
	if err := s.authRepo.UpdateResetToken(email, resetToken, expiry); err != nil {
		return "", errors.New(constants.ErrResetTokenFailed)
	}

	return resetToken, nil
}

// ResetPassword handles password reset using a token
func (s *AuthServiceImpl) ResetPassword(token, newPassword string) error {
	// Check if password is strong
	if !authPkg.IsStrongPassword(newPassword) {
		return errors.New(constants.ErrStrongPassword)
	}

	// Find the user by reset token
	user, err := s.authRepo.GetUserByResetToken(token)
	if err != nil {
		return errors.New(constants.ErrInvalidToken)
	}

	// Check if token is expired
	if user.ResetPasswordExpires == nil || user.ResetPasswordExpires.Before(time.Now()) {
		return errors.New(constants.ErrTokenExpired)
	}

	// Hash the new password
	hashedPassword, err := authPkg.HashPassword(newPassword)
	if err != nil {
		return errors.New(constants.ErrPasswordHashFailed)
	}

	// Update user with new password
	if err := s.authRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return errors.New(constants.ErrDatabase)
	}

	return nil
}
