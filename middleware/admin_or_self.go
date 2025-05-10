package middleware

import (
	"go-gin-starter/models"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminOrSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
			c.Abort()
			return
		}

		currentUserID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
			c.Abort()
			return
		}

		user, err := services.GetUserByID(currentUserID)
		if err != nil {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
			c.Abort()
			return
		}

		// Get target user ID from URL param
		targetIDParam := c.Param("id")
		targetUserID, err := uuid.Parse(targetIDParam)
		if err != nil {
			httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
			c.Abort()
			return
		}

		// Check if current user is admin or self
		if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin && currentUserID != targetUserID {
			httpPkg.RespondError(c, http.StatusForbidden, constants.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
