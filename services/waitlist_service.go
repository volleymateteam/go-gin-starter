package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/models"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"

	"github.com/google/uuid"
)

// SubmitWaitlistEntry handles inserting a new waitlist entry
func SubmitWaitlistEntry(email, source string) error {
	exists, err := repositories.IsEmailAlreadyInWaitlist(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(utils.ErrAlreadyInWaitlist)
	}
	return repositories.CreateWaitlistEntry(email, source)
}

// GetAllWaitlistEntries returns all waitlist entries mapped to DTO
func GetAllWaitlistEntries() ([]dto.WaitlistEntryResponse, error) {
	entries, err := repositories.GetAllWaitlistEntries()
	if err != nil {
		return nil, err
	}

	var responses []dto.WaitlistEntryResponse
	for _, entry := range entries {
		responses = append(responses, dto.WaitlistEntryResponse{
			ID:     entry.ID.String(),
			Email:  entry.Email,
			Source: entry.Source,
		})
	}
	return responses, nil
}

// ApproveWaitlistEntry creates a user and removes the waitlist entry
func ApproveWaitlistEntry(id uuid.UUID) error {
	entry, err := repositories.FindWaitlistEntryByID(id)
	if err != nil {
		return err
	}

	// Create user with email as username and random password
	randomPassword := utils.GenerateRandomPassword()
	hashedPassword, err := utils.HashPassword(randomPassword)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    entry.Email,
		Username: entry.Email, // Can change later by user
		Password: hashedPassword,
		Role:     models.RolePlayer,
		Gender:   models.GenderOther,
	}

	if err := repositories.CreateUser(&user); err != nil {
		return err
	}

	// Remove waitlist entry
	return repositories.DeleteWaitlistEntryByID(id)
}

// RejectWaitlistEntry removes the waitlist entry
func RejectWaitlistEntry(id uuid.UUID) error {
	return repositories.DeleteWaitlistEntryByID(id)
}
