package request

type CreateRoom struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password"`
}

type LoginRoom struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password"`
}

type UpdateRoom struct {
	CreateRoom
}

type ListRoom struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDRoom struct {
	ID string `uri:"id" binding:"required"`
}
