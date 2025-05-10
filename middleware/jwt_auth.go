package middleware

import (
	"net/http"
	"strings"

	authPkg "go-gin-starter/pkg/auth"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authPkg.ParseJWT(tokenStr)
		if err != nil {
			httpPkg.RespondError(c, http.StatusUnauthorized, constants.ErrInvalidToken)
			c.Abort()
			return
		}

		// Store UUID user ID in context directly
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
