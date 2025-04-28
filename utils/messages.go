package utils

// Common error messages
const (
	ErrUnauthorized          = "Unauthorized"
	ErrForbidden             = "Access denied"
	ErrInvalidCredentials    = "Invalid email or password"
	ErrUserNotFound          = "User not found"
	ErrInvalidUserID         = "Invalid user ID"
	ErrInvalidGender         = "Invalid gender value"
	ErrInvalidRole           = "Invalid role value"
	ErrInternalServer        = "Internal Server Error"
	ErrInvalidToken          = "Invalid or expired token"
	ErrTokenExpired          = "Reset token has expired"
	ErrTokenGenerationFailed = "Token generation failed"
	ErrResetTokenFailed      = "Failed to generate reset token"
	ErrInvalidInput          = "Invalid input"
	ErrAvatarTooLarge        = "Avatar file is too large. Max size is 2MB"
	ErrInvalidFileType       = "Only .jpg, .jpeg, and .png formats are allowed"
	ErrUploadFailed          = "Failed to upload file"
	ErrDatabase              = "Database error"
	ErrPasswordMismatch      = "Invalid old password"
	ErrStrongPassword        = "Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters."
	ErrAlreadyInWaitlist     = "already in waitlist"
)

// Common success messages
const (
	MsgProfileFetched    = "Profile fetched successfully"
	MsgUserRegistered    = "User registered successfully"
	MsgUserUpdated       = "User updated successfully"
	MsgUserDeleted       = "User deleted successfully"
	MsgPasswordChanged   = "Password changed successfully"
	MsgAvatarUploaded    = "Avatar uploaded successfully"
	MsgResetTokenCreated = "Reset token generated successfully"
	MsgPasswordReset     = "Password reset successfully"
	MsgUsersFetched      = "Users fetched successfully"
	MsgWaitlistFetched   = "Waitlist fetched successfully"
	MsgWaitlistSubmitted = "Waitlist submitted successfully"
)
