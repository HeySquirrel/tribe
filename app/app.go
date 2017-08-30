package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

type App struct {
	Gui  *gocui.Gui
	Done chan struct{}
}

func New() *App {
	a := new(App)
	a.Done = make(chan struct{})

	var err error
	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.Gui.SetManagerFunc(a.Layout)
	a.setKeyBindings()

	return a
}

func (a *App) Loop() {
	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) Close() {
	close(a.Done)
	a.Gui.Close()
}

func (a *App) updateCurrentFile() {
	file := a.currentFileSelection()
	a.updateLogs(file)
}
