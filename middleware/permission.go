package middleware

import (
	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
			c.Abort()
			return
		}

		user, err := services.GetUserByID(userID)
		if err != nil {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
			c.Abort()
			return
		}

		if !authPkg.HasPermission(user, permission) {
			httpPkg.RespondError(c, http.StatusForbidden, constants.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
