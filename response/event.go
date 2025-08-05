package response

import (
	"time"
)

type EventResponse struct {
	ID          uint64     `json:"taskId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Complete    *bool      `json:"complete"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updateAt"`
	Location    string     `json:"location"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     time.Time  `json:"endTime"`
}
