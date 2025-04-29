package utils

// Common error messages
const (
	ErrUnauthorized          = "unauthorized"
	ErrForbidden             = "access denied"
	ErrInvalidCredentials    = "invalid email or password"
	ErrUserNotFound          = "user not found"
	ErrInvalidUserID         = "invalid user ID"
	ErrInvalidGender         = "invalid gender value"
	ErrInvalidRole           = "invalid role value"
	ErrInternalServer        = "internal Server Error"
	ErrInvalidToken          = "invalid or expired token"
	ErrTokenExpired          = "reset token has expired"
	ErrTokenGenerationFailed = "token generation failed"
	ErrResetTokenFailed      = "failed to generate reset token"
	ErrInvalidInput          = "invalid input"
	ErrAvatarTooLarge        = "avatar file is too large. Max size is 2MB"
	ErrInvalidFileType       = "only .jpg, .jpeg, and .png formats are allowed"
	ErrUploadFailed          = "failed to upload file"
	ErrDatabase              = "database error"
	ErrPasswordMismatch      = "invalid old password"
	ErrStrongPassword        = "password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters."
	ErrAlreadyInWaitlist     = "already in waitlist"
	ErrSeasonNotFound        = "season not found"
)

// Common success messages
const (
	MsgProfileFetched    = "profile fetched successfully"
	MsgUserRegistered    = "user registered successfully"
	MsgUserUpdated       = "user updated successfully"
	MsgUserDeleted       = "user deleted successfully"
	MsgPasswordChanged   = "password changed successfully"
	MsgAvatarUploaded    = "avatar uploaded successfully"
	MsgResetTokenCreated = "reset token generated successfully"
	MsgPasswordReset     = "password reset successfully"
	MsgUsersFetched      = "users fetched successfully"
	MsgWaitlistFetched   = "waitlist fetched successfully"
	MsgWaitlistSubmitted = "waitlist submitted successfully"
	MsgWaitlistApproved  = "waitlist approved successfully"
	MsgWaitlistRejected  = "waitlist rejected successfully"
	MsgSeasonCreated     = "season created successfully"
	MsgSeasonUpdated     = "season updated successfully"
	MsgSeasonDeleted     = "season deleted successfully"
	MsgSeasonFetched     = "season fetched successfully"
	MsgSeasonsFetched    = "seasons fetched successfully"
)
