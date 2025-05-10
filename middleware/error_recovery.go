package middleware

import (
	"log"
	"net/http"

	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"

	"github.com/gin-gonic/gin"
)

// ErrorRecovery catches all panics and returns a JSON 500
func ErrorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] %v\n", rec)
				httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrInternalServer)
			}
		}()
		c.Next()
	}
}
