package context

import (
	"prick/internal/prick"
	"prick/internal/prick/config"
	"prick/internal/prick/interfaces"
)

type BubbleContext struct {
	Config           *config.Config
	User             string
	ScreenHeight     int
	ScreenWidth      int
	ContentHeight    int
	Api              *prick.Api
	SelectedResource interfaces.Prickable
}
