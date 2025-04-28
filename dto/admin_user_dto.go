package dto

import (
	"go-gin-starter/models"

	"github.com/google/uuid"
)

type AdminUserResponse struct {
	ID        uuid.UUID         `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Gender    models.GenderEnum `json:"gender"`
	Role      models.RoleEnum   `json:"role"`
	AvatarURL string            `json:"avatar_url"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	DeletedAt string            `json:"deleted_at,omitempty"`
}

type AdminUpdateUserInput struct {
	Username string            `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string            `json:"email" binding:"omitempty,email"`
	Gender   models.GenderEnum `json:"gender" binding:"omitempty"`
	Role     models.RoleEnum   `json:"role" binding:"omitempty"`
}
