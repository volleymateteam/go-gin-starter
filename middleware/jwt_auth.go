package middleware

import (
	"net/http"
	"strings"

	"go-gin-starter/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, utils.ErrInvalidToken)
			c.Abort()
			return
		}

		// Store UUID user ID in context directly
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
