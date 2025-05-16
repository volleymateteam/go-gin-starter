package services

import (
	"errors"
	"go-gin-starter/dto"
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	"go-gin-starter/repositories"

	"github.com/google/uuid"
)

// WaitlistService defines the interface for waitlist-related business logic
type WaitlistService interface {
	SubmitWaitlistEntry(email, source string) error
	GetAllWaitlistEntries() ([]dto.WaitlistEntryResponse, error)
	ApproveWaitlistEntry(id uuid.UUID) error
	RejectWaitlistEntry(id uuid.UUID) error
}

// WaitlistServiceImpl implements WaitlistService
type WaitlistServiceImpl struct {
	waitlistRepo repositories.WaitlistRepository
	userService  UserService
}

// NewWaitlistService creates a new instance of WaitlistService
func NewWaitlistService(waitlistRepo repositories.WaitlistRepository, userService UserService) WaitlistService {
	return &WaitlistServiceImpl{
		waitlistRepo: waitlistRepo,
		userService:  userService,
	}
}

// SubmitWaitlistEntry handles inserting a new waitlist entry
func (s *WaitlistServiceImpl) SubmitWaitlistEntry(email, source string) error {
	exists, err := s.waitlistRepo.IsEmailAlreadyInWaitlist(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(constants.ErrAlreadyInWaitlist)
	}
	return s.waitlistRepo.CreateWaitlistEntry(email, source)
}

// GetAllWaitlistEntries returns all waitlist entries mapped to DTO
func (s *WaitlistServiceImpl) GetAllWaitlistEntries() ([]dto.WaitlistEntryResponse, error) {
	entries, err := s.waitlistRepo.GetAllWaitlistEntries()
	if err != nil {
		return nil, err
	}

	// Always initialize an empty slice, not nil
	responses := make([]dto.WaitlistEntryResponse, 0, len(entries))

	for _, entry := range entries {
		responses = append(responses, dto.WaitlistEntryResponse{
			ID:        entry.ID.String(),
			Email:     entry.Email,
			Source:    entry.Source,
			CreatedAt: entry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// ApproveWaitlistEntry creates a user and removes the waitlist entry
func (s *WaitlistServiceImpl) ApproveWaitlistEntry(id uuid.UUID) error {
	entry, err := s.waitlistRepo.FindWaitlistEntryByID(id)
	if err != nil {
		return err
	}

	// Create user with email as username and random password
	randomPassword := authPkg.GenerateRandomPassword()

	// Use the user service to create the user
	_, err = s.userService.CreateUser(entry.Email, entry.Email, randomPassword, "other")
	if err != nil {
		return err
	}

	// Remove waitlist entry
	return s.waitlistRepo.DeleteWaitlistEntryByID(id)
}

// RejectWaitlistEntry removes the waitlist entry
func (s *WaitlistServiceImpl) RejectWaitlistEntry(id uuid.UUID) error {
	return s.waitlistRepo.DeleteWaitlistEntryByID(id)
}
