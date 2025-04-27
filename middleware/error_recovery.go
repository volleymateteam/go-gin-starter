package middleware

import (
	"log"
	"net/http"

	"go-gin-starter/utils"

	"github.com/gin-gonic/gin"
)

// ErrorRecovery catches all panics and returns a JSON 500
func ErrorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] %v\n", rec)
				utils.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
