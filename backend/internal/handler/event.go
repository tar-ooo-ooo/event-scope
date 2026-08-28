package handler

import (
	"event-scope/backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEventHandler(context *gin.Context) {
	var req model.CreateEventRequest

	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 發佈到 Kafka

	context.JSON(http.StatusAccepted, gin.H{"event_id": req.EventId})
}