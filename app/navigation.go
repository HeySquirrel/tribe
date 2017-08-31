package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) NextFile(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}

		a.currentFileChanged()
	}
	return nil
}

func (a *App) PreviousFile(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}

		a.currentFileChanged()
	}
	return nil
}

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
