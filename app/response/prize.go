package response

type ListPrize struct {
	ID       string `bun:"id" json:"id"`
	Name     string `bun:"name" json:"name"`
	ImageURL string `bun:"image_url" json:"image_url"`
	Quantity int64  `bun:"quantity" json:"quantity"`
	RoomID   string `bun:"room_id" json:"room_id"`
}

type Prize struct {
	ID        string `bun:"id" json:"id"`
	Name      string `bun:"name" json:"name"`
	ImageURL  string `bun:"image_url" json:"image_url"`
	Quantity  int64  `bun:"quantity" json:"quantity"`
	RoomID    string `bun:"room_id" json:"room_id"`
	CreatedAt string `bun:"created_at" json:"created_at"`
}
