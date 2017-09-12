package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
