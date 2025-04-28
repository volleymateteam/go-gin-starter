package middleware

import (
	"go-gin-starter/models"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminOrSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		currentUserID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			utils.RespondError(c, http.StatusInternalServerError, utils.ErrInvalidUserID)
			c.Abort()
			return
		}

		user, err := services.GetUserByID(currentUserID)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		// Get target user ID from URL param
		targetIDParam := c.Param("id")
		targetUserID, err := uuid.Parse(targetIDParam)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
			c.Abort()
			return
		}

		// Check if current user is admin or self
		if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin && currentUserID != targetUserID {
			utils.RespondError(c, http.StatusForbidden, utils.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
