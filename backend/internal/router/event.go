package router

import (
	"github.com/gin-gonic/gin"

	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/sse"
)

func eventRoutes(rg *gin.RouterGroup, broker *sse.Broker) {
	rg.POST("", handler.CreateEventHandler).GET("/stream", handler.StreamEventHandler(broker))

}
