package router

import (
	"event-scope/backend/internal/config"
	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/middleware"
	"event-scope/backend/internal/sse"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors(cfg))

	router.GET("/healthz", func(context *gin.Context) {
		handler.Res(context, http.StatusOK, "API is healthy", nil)
	})

	// 建立 SSE Broker
	broker := sse.NewBroker()

	eventRoutes(router.Group("/event"), broker)

	return router
}
