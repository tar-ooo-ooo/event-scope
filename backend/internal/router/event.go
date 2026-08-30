package router

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/sse"
)

func eventRoutes(rg *gin.RouterGroup, sseB *sse.Broker, wr *kafka.Writer) {
	rg.POST("", handler.CreateEventHandler(sseB, wr)).GET("/stream", handler.StreamEventHandler(sseB))
}
