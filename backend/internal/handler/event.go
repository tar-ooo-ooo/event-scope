package handler

import (
	"event-scope/backend/internal/model"
	"event-scope/backend/internal/sse"
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

	Res(context, http.StatusAccepted, gin.H{"event_id": req.EventId}, nil)
}

func StreamEventHandler(broker *sse.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 設定 SSE 標頭
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		client := broker.Subscribe()
		defer broker.Unsubscribe(client)

		for {
			select {
			case event := <-client: // 從 broker 的 client channel 接收事件
				c.SSEvent("message", event)

				// 立即將事件發送給前端
				c.Writer.Flush()
			case <-c.Request.Context().Done(): // 前端斷線時，退出循環
				return
			}
		}
	}
}
