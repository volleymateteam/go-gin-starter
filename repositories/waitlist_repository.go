package repositories

import (
	"go-gin-starter/database"
	"go-gin-starter/models"

	"github.com/google/uuid"
)

// WaitlistRepository defines the interface for waitlist data operations
type WaitlistRepository interface {
	IsEmailAlreadyInWaitlist(email string) (bool, error)
	CreateWaitlistEntry(email, source string) error
	GetAllWaitlistEntries() ([]models.WaitlistEntry, error)
	FindWaitlistEntryByID(id uuid.UUID) (*models.WaitlistEntry, error)
	DeleteWaitlistEntryByID(id uuid.UUID) error
}

// GormWaitlistRepository implements WaitlistRepository using GORM
type GormWaitlistRepository struct{}

// NewWaitlistRepository creates a new instance of WaitlistRepository
func NewWaitlistRepository() WaitlistRepository {
	return &GormWaitlistRepository{}
}

// CreateWaitlistEntry inserts a new waitlist record
func (r *GormWaitlistRepository) CreateWaitlistEntry(email, source string) error {
	entry := models.WaitlistEntry{
		ID:     uuid.New(),
		Email:  email,
		Source: source,
	}
	return database.DB.Create(&entry).Error
}

// IsEmailAlreadyInWaitlist checks if email already exists
func (r *GormWaitlistRepository) IsEmailAlreadyInWaitlist(email string) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.WaitlistEntry{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllWaitlistEntries retrieves all waitlist records (ID, Email, Source)
func (r *GormWaitlistRepository) GetAllWaitlistEntries() ([]models.WaitlistEntry, error) {
	var entries []models.WaitlistEntry
	if err := database.DB.Order("created_at desc").Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// FindWaitlistEntryByID finds a waitlist entry by ID
func (r *GormWaitlistRepository) FindWaitlistEntryByID(id uuid.UUID) (*models.WaitlistEntry, error) {
	var entry models.WaitlistEntry
	if err := database.DB.First(&entry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

// DeleteWaitlistEntryByID deletes a waitlist entry by ID
func (r *GormWaitlistRepository) DeleteWaitlistEntryByID(id uuid.UUID) error {
	return database.DB.Delete(&models.WaitlistEntry{}, "id = ?", id).Error
}
