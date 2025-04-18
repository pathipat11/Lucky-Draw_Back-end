package model

import (
	"github.com/uptrace/bun"
)

type Prize struct {
	bun.BaseModel `bun:"table:prizes"`

	ID       string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name     string `bun:"name,notnull"`
	ImageURL string `bun:"image_url,notnull"`
	Quantity int64  `bun:"quantity,notnull"`
	RoomID   string `bun:"room_id,notnull"`

	Room *Room `bun:"rel:belongs-to,join:room_id=id"`

	Winners        []Winner        `bun:"rel:has-many,join:id=prize_id"`
	DrawConditions []DrawCondition `bun:"rel:has-many,join:id=prize_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
