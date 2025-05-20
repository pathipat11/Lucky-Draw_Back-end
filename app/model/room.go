package model

import (
	"github.com/uptrace/bun"
)

type Room struct {
	bun.BaseModel `bun:"table:rooms"`

	ID       string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name     string `bun:"name,notnull"`
	Password string `bun:"password,nullzero" json:"-"`

	Players        []Player        `bun:"rel:has-many,join:id=room_id"`
	Prizes         []Prize         `bun:"rel:has-many,join:id=room_id"`
	DrawConditions []DrawCondition `bun:"rel:has-many,join:id=room_id"`
	Winners        []Winner        `bun:"rel:has-many,join:id=room_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
