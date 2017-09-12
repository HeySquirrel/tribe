package app

import (
	"github.com/jroimartin/gocui"
)

const (
	changesView = "changes"
)

func (a *App) Layout(g *gocui.Gui) error {
	return a.setKeyBindings()
}
