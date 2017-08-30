package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

type App struct {
	Gui *gocui.Gui
}

func New() *App {
	a := new(App)

	var err error
	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.Gui.SetManagerFunc(a.Layout)

	return a
}

func (a *App) Loop() {
	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) Close() {
	a.Gui.Close()
}
