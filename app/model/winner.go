package model

import (
	"github.com/uptrace/bun"
)

type Winner struct {
	bun.BaseModel `bun:"table:winners"`

	ID              string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	RoomID          string `bun:"room_id, notnull"`
	PlayerID        string `bun:"player_id, notnull"`
	PrizeID         string `bun:"prize_id, notnull"`
	DrawConditionID string `bun:"draw_condition_id, notnull"`

	Room          *Room          `bun:"rel:belongs-to,join:room_id=id"`
	Player        *Player        `bun:"rel:belongs-to,join:player_id=id"`
	Prize         *Prize         `bun:"rel:belongs-to,join:prize_id=id"`
	DrawCondition *DrawCondition `bun:"rel:belongs-to,join:draw_condition_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
