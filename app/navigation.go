package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) ShowDebug(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop(debugView)
	g.SetCurrentView(debugView)
	return nil
}

func (a *App) HideDebug(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom(debugView)
	g.SetCurrentView(changesView)
	return nil
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
