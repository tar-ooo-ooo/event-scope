package model

import "time"

type APIResponse struct {
	Data      any       `json:"data"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}
