package model

type CreateEventRequest struct {
	EventId string `json:"event_id" binding:"required"`
	Type string `json:"type" binding:"required"`
}