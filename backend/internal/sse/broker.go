package sse

import (
	"event-scope/backend/internal/model"
	"sync"
)

type Broker struct {
	// mutex，多人同時連線、斷線或廣播時，用它保護 clients，避免同時修改 map 造成 crash
	mu sync.RWMutex

	// 每個 SSE 前端連線對應一條 channel，當有新的事件時，會透過這個 channel 廣播給前端
	clients map[chan model.EventResult]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan model.EventResult]struct{}),
	}
}

func (b *Broker) Subscribe() chan model.EventResult {
	client := make(chan model.EventResult, 1)

	b.mu.Lock()                    // 加鎖，避免多個 client 同時修改 clients map
	b.clients[client] = struct{}{} // 將新的 client channel 加入 clients map
	b.mu.Unlock()                  // 解鎖

	return client
}

func (b *Broker) Unsubscribe(client chan model.EventResult) {
	b.mu.Lock()
	delete(b.clients, client)
	b.mu.Unlock()
}

func (b *Broker) Publish(event model.EventResult) {
	b.mu.RLock()         // 讀鎖，避免在廣播時有 client 被加入或移除，用 RLock 是因為不需要修改 clients map，只是讀取它
	defer b.mu.RUnlock() // 函式結束時自動解除讀鎖

	for client := range b.clients {
		select {
		case client <- event: // 將事件發送到 client channel
		default:
			// 如果 client channel 已滿，則跳過該 client，避免阻塞
		}
	}
}
