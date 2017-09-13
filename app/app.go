package app

import (
	"github.com/heysquirrel/tribe/git"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/widgets"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"time"
)

type App struct {
	Gui       *gocui.Gui
	Done      chan struct{}
	Log       *tlog.Log
	Git       *git.Repo
	Changes   *widgets.ChangesView
	DebugView *widgets.DebugView
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

	a.Changes = widgets.NewChangesView(a.Gui)
	a.DebugView = widgets.NewDebugView(a.Gui)

	associatedFiles := widgets.NewAssociatedFilesView(a.Gui)
	recentContributors := widgets.NewRecentContributorsView(a.Gui)
	relatedWork := widgets.NewRelatedWorkView(a.Gui)
	logs := widgets.NewLogsView(a.Gui)
	legend := widgets.NewLegendView(a.Gui)
	feed := widgets.NewFeedView(a.Gui)

	a.Changes.AddListener(func(selectedFile *git.File) {
		associatedFiles.UpdateRelatedFiles(selectedFile.Related)
		recentContributors.UpdateContributors(selectedFile.Contributors)
		relatedWork.UpdateRelatedWork(selectedFile.WorkItems)
	})

	a.Gui.SetManager(
		a.Changes,
		associatedFiles,
		recentContributors,
		relatedWork,
		logs,
		legend,
		feed,
		a.DebugView,
	)

	a.setKeyBindings()

	return a
}

func (a *App) Debug(message string) {
	a.Log.Add(message)
	a.DebugView.UpdateDebug(a.Log.Entries())
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
	a.Changes.SetChanges(a.Git.Changes())
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			a.Debug("Checking for changes")
			a.Changes.SetChanges(a.Git.Changes())
		}
	}

}
