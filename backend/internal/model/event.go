package model

type CreateEventRequest struct {
	EventId string `json:"event_id" binding:"required"`
}

type EventResult struct {
	EventId string `json:"event_id"`
	Success bool   `json:"success"`
}
