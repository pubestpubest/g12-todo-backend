package request

type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1" json:"page"`
	Limit int `form:"limit" binding:"min=1,max=100" json:"limit"`
}
