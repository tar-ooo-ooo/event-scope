package middleware

import (
	"net/http"

	"event-scope/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func Cors(cfg config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.GetHeader("Origin") == cfg.FrontendEndpoint {
			c.Header("Access-Control-Allow-Origin", cfg.FrontendEndpoint)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")

			if c.Request.Method == http.MethodOptions {
				c.Status(http.StatusNoContent)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
