package controller

import (
	"app/app/controller/draw_condition"
	"app/app/controller/player"
	"app/app/controller/prize"
	"app/app/controller/room"
	"app/app/controller/winner"
	"app/config"
)

type Controller struct {
	RoomCtl          *room.Controller
	PlayerCtl        *player.Controller
	PrizeCtl         *prize.Controller
	DrawConditionCtl *draw_condition.Controller
	WinnerCtl        *winner.Controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{
		RoomCtl:          room.NewController(db),
		PlayerCtl:        player.NewController(db),
		PrizeCtl:         prize.NewController(db),
		DrawConditionCtl: draw_condition.NewController(db),
		WinnerCtl:        winner.NewController(db),
		// Other controllers...
	}
}
