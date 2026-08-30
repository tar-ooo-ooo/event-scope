package router

import (
	"net/http"

	"event-scope/backend/internal/config"
	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/middleware"
	"event-scope/backend/internal/sse"

	"github.com/gin-gonic/gin"
	kafkago "github.com/segmentio/kafka-go"
)

func Setup(cfg config.Config, wr *kafkago.Writer) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestLogger(),
		gin.Recovery(),
		middleware.Cors(cfg),
	)

	router.GET("/healthz", func(ctx *gin.Context) {
		handler.Res(ctx, http.StatusOK, "API is healthy", nil)
	})

	// 建立 SSE Broker
	sseBroker := sse.NewBroker()

	eventRoutes(router.Group("/event"), sseBroker, wr)

	return router
}
