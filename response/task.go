package response

import "time"

type TaskResponse struct {
	ID          uint64     `json:"taskId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      bool       `json:"status"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updateAt"`
}
