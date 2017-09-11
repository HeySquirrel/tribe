package app

import (
	"github.com/heysquirrel/tribe/git"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"time"
)

type App struct {
	Gui  *gocui.Gui
	Done chan struct{}
	Log  *tlog.Log
	Git  *git.Repo
}

func New() *App {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	a := new(App)
	a.Done = make(chan struct{})
	a.Log = tlog.New()

	a.Git, err = git.New(pwd, a.Log)
	if err != nil {
		log.Panicln(err)
	}

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	a.Gui.SetManager(a)

	return a
}

func (a *App) Debug(message string) {
	a.Log.Add(message)
	a.UpdateDebug(a.Log.Entries())
}

func (a *App) Loop() {
	go a.checkForChanges()

	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) Close() {
	close(a.Done)
	a.Gui.Close()
}

func (a *App) checkForChanges() {
	a.UpdateChanges(a.Git.Changes())
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			a.Debug("Checking for changes")
			a.UpdateChanges(a.Git.Changes())
		}
	}

}

func (a *App) currentFileChanged() {
	file := a.currentFileSelection()

	go func(app *App, file string) {
		files, workItems, contributors := app.Git.Related(file)
		app.UpdateContributors2(contributors)
		app.UpdateRelatedFiles(files)
		app.UpdateRelatedWork(workItems)
	}(a, file)
}
