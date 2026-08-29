package handler

import (
	"event-scope/backend/internal/model"
	"event-scope/backend/internal/sse"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func CreateEventHandler(b *sse.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateEventRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			Res(c, 400, nil, err)
			return
		}

		success := rand.Intn(100) < 75

		result := model.EventResult{
			EventId: req.EventId,
			Success: success,
		}
		b.Publish(result)

		Res(c, 200, "Event created and published", nil)
	}
}

func StreamEventHandler(b *sse.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 設定 SSE 標頭
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		client := b.Subscribe()
		defer b.Unsubscribe(client)

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
