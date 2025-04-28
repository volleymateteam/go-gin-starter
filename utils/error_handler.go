package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse defines the format we want for all errors
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SuccessResponse defines a standard success format
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// RespondError sends a standardized error response
func RespondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// RespondSuccess sends a standardized success response
func RespondSuccess(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(code, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
