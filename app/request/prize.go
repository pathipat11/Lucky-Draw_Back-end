package request

type CreatePrize struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Quantity int64  `json:"quantity"`
	RoomID   string `json:"room_id"`
}

type UpdatePrize struct {
	CreatePrize
}

type ListPrize struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDPrize struct {
	ID string `uri:"id" binding:"required"`
}
