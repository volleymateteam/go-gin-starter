package services

import (
	"errors"
	"go-gin-starter/repositories"
	"go-gin-starter/utils"
)

// SubmitWaitlistEntry creates a new waitlist entry
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

// GetAllWaitlistEntries retrieves all waitlist entries
func GetAllWaitlistEntries() ([]string, error) {
	entries, err := repositories.GetAllWaitlistEmails()
	if err != nil {
		return nil, err
	}
	return entries, nil
}
