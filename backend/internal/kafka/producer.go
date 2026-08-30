package kafka

import (
	"context"
	"encoding/json"
	"event-scope/backend/internal/model"

	kafkago "github.com/segmentio/kafka-go"
)

func NewWriter(b string, t string) *kafkago.Writer {
	return &kafkago.Writer{
		Addr:         kafkago.TCP(b),
		Topic:        t,
		RequiredAcks: kafkago.RequireAll, // Kafka 會把訊息儲存到磁碟的 topic 裡，等到所有 partition 的 replica 都收到訊息後，才會回覆給 client，確保訊息不會遺失
		// BatchSize: 100, // 批次寫入訊息的大小，預設是 100，若要即時寫入訊息，可以設定為 1
	}
}

func Publish(ctx context.Context, wr *kafkago.Writer, e model.CreateEventRequest) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return wr.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(e.EventId),
		Value: payload,
	})
}
