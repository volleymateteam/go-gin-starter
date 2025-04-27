package middleware

import (
	"go-gin-starter/models"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			utils.RespondError(c, http.StatusInternalServerError, "Invalid user ID")
			c.Abort()
			return
		}

		user, err := services.GetUserByID(userID)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin {
			utils.RespondError(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
