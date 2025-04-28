package services

import (
	"errors"
	"go-gin-starter/dto"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"
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
