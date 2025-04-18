package response

type ListWinner struct {
	ID              string `bun:"id" json:"id"`
	RoomID          string `bun:"room_id" json:"room_id"`
	PlayerID        string `bun:"player_id" json:"player_id"`
	PrizeID         string `bun:"prize_id" json:"prize_id"`
	DrawConditionID string `bun:"draw_condition_id" json:"draw_condition_id"`
}
