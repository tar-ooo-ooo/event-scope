package router

import (
	"github.com/gin-gonic/gin"

	"event-scope/backend/internal/handler"
	"event-scope/backend/internal/sse"
)

func eventRoutes(rg *gin.RouterGroup, b *sse.Broker) {
	rg.POST("", handler.CreateEventHandler(b)).GET("/stream", handler.StreamEventHandler(b))

}
