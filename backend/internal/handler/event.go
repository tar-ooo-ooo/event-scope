package handler

import (
	"event-scope/backend/internal/model"
	"event-scope/backend/internal/sse"
	"net/http"

	kafka "event-scope/backend/internal/kafka"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/gin-gonic/gin"
)

func CreateEventHandler(b *sse.Broker, wr *kafkago.Writer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.CreateEventRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			Res(ctx, http.StatusBadRequest, nil, err)
			return
		}

		if err := kafka.Publish(ctx, wr, req); err != nil {
			Res(ctx, http.StatusServiceUnavailable, nil, err)
			return
		}

		result := model.EventResult{
			EventId: req.EventId,
			Success: true,
		}

		b.Publish(result)

		Res(ctx, http.StatusAccepted, "Event created and published", nil)
	}
}

func StreamEventHandler(b *sse.Broker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 設定 SSE 標頭
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")

		client := b.Subscribe()
		defer b.Unsubscribe(client)

		for {
			select {
			case event := <-client: // 從 broker 的 client channel 接收事件
				ctx.SSEvent("message", event)

				// 立即將事件發送給前端
				ctx.Writer.Flush()
			case <-ctx.Request.Context().Done(): // 前端斷線時，退出循環
				return
			}
		}
	}
}
