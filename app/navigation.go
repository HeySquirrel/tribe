package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) ShowDebug(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("debug")
	g.SetCurrentView("debug")
	return nil
}

func (a *App) HideDebug(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("debug")
	g.SetCurrentView(changesView)
	return nil
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
