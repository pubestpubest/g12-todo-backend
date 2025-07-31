package response

type Response[T any] struct {
	Status  string `json:"status"`  // "success" or "error"
	Message string `json:"message"` // human-readable message
	Data    T      `json:"data"`    // generic payload
}

type PaginatedResponse[T any] struct {
	Status     string     `json:"status"`     // "success" or "error"
	Message    string     `json:"message"`    // human-readable message
	Data       []T        `json:"data"`       // array of payloads
	Pagination Pagination `json:"pagination"` // pagination info
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
