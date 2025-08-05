package request

import (
	"time"
)

type EventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location" binding:"required"`
	StartTime   time.Time `json:"startTime" binding:"required"`
	EndTime     time.Time `json:"endTime" binding:"required"`
	Complete    *bool     `json:"complete" binding:"required"`
}
