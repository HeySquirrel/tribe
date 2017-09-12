package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

func (a *App) setKeyBindings() error {
	err := a.Gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, a.ShowDebug)
	if err != nil {
		log.Panicln(err)
	}

	err = a.Gui.SetKeybinding("debug", gocui.KeyF1, gocui.ModNone, a.HideDebug)
	if err != nil {
		log.Panicln(err)
	}

	err = a.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, a.quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}
