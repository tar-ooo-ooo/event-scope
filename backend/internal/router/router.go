package router

import (
	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/sse"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", func(context *gin.Context) {
		handler.Res(context, http.StatusOK, "API is healthy", nil)
	})

	// 建立 SSE Broker
	broker := sse.NewBroker()

	eventRoutes(router.Group("/event"), broker)

	return router
}
