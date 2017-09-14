package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

func (a *App) setKeyBindings() error {
	// show := func(g *gocui.Gui, v *gocui.View) error { a.DebugView.Show(); return nil }
	// err := a.Gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, show)
	// if err != nil {
	// 	log.Panicln(err)
	// }

	quit := func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }
	err := a.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}
