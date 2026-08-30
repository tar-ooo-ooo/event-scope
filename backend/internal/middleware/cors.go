package middleware

import (
	"net/http"

	"event-scope/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func Cors(cfg config.Config) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		if ctx.GetHeader("Origin") == cfg.FrontendEndpoint {
			ctx.Header("Access-Control-Allow-Origin", cfg.FrontendEndpoint)
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type")

			if ctx.Request.Method == http.MethodOptions {
				ctx.Status(http.StatusNoContent)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
