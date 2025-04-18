package response

type ListPlayer struct {
	ID        string `bun:"id" json:"id"`
	Prefix    string `bun:"prefix" json:"prefix"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	MemberID  string `bun:"member_id" json:"member_id"`
	Position  string `bun:"position" json:"position"`
	RoomID    string `bun:"room_id" json:"room_id"`
	IsActive  bool   `bun:"is_active" json:"is_active"`
}
