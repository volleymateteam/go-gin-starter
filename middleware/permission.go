package middleware

import (
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			utils.RespondError(c, http.StatusInternalServerError, utils.ErrInvalidUserID)
			c.Abort()
			return
		}

		user, err := services.GetUserByID(userID)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		if !utils.HasPermission(user, permission) {
			utils.RespondError(c, http.StatusForbidden, utils.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
