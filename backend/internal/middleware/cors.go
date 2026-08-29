package middleware

import (
	"event-scope/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func Cors(cfg config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.GetHeader("Origin") == cfg.FrontendEndpoint {
			c.Header("Access-Control-Allow-Origin", cfg.FrontendEndpoint)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
		}
	}
}
