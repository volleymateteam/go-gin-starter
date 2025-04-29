package dto

type CreateWaitlistEntryInput struct {
	Email  string `json:"email" binding:"required,email"`
	Source string `json:"source" binding:"omitempty"` // Optional field
}

type WaitlistEntryResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Source string `json:"source,omitempty"` // Optional field
}
