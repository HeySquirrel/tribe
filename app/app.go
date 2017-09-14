package app

import (
	"github.com/heysquirrel/tribe/app/widgets"
	"github.com/heysquirrel/tribe/blame"
	"github.com/jroimartin/gocui"
	"log"
)

type App struct {
	Gui  *gocui.Gui
	Done chan struct{}
	// Log       *tlog.Log
	// Git       *git.Repo
	// Changes   *widgets.ChangesView
	// DebugView *widgets.DebugView
}

func New(filename string) *App {
	// pwd, err := os.Getwd()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// usr, err := user.Current()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// configFile := filepath.Join(usr.HomeDir, ".tribe")
	// config, err := ioutil.ReadFile(configFile)
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// api := rally.New(string(config))

	a := new(App)
	a.Done = make(chan struct{})
	// a.Log = tlog.New()

	// a.Git, err = git.New(pwd, a.Log, api)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	blame, err := blame.New(filename)
	if err != nil {
		log.Panicln(err)
	}

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	// a.Changes = widgets.NewChangesView(a.Gui)
	// a.DebugView = widgets.NewDebugView(a.Gui)

	// associatedFiles := widgets.NewAssociatedFilesView(a.Gui)
	// recentContributors := widgets.NewRecentContributorsView(a.Gui)
	// relatedWork := widgets.NewRelatedWorkView(a.Gui)
	// logs := widgets.NewLogsView(a.Gui)
	// legend := widgets.NewLegendView(a.Gui)
	// feed := widgets.NewFeedView(a.Gui)

	// a.Changes.AddListener(func(selectedFile *git.File) {
	// 	associatedFiles.UpdateRelatedFiles(selectedFile.Related)
	// 	recentContributors.UpdateContributors(selectedFile.Contributors)
	// 	relatedWork.UpdateRelatedWork(selectedFile.WorkItems)
	// })

	source := widgets.NewSourceCodeView(a.Gui, blame)

	a.Gui.SetManager(
		source,
		// a.Changes,
		// associatedFiles,
		// recentContributors,
		// relatedWork,
		// logs,
		// legend,
		// feed,
		// a.DebugView,
	)

	a.setKeyBindings()

	return a
}

func (a *App) Debug(message string) {
	// a.Log.Add(message)
	// a.DebugView.UpdateDebug(a.Log.Entries())
}

func (a *App) Loop() {
	// go a.checkForChanges()

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
	// a.Changes.SetChanges(a.Git.Changes())
	// for {
	// 	select {
	// 	case <-a.Done:
	// 		return
	// 	case <-time.After(10 * time.Second):
	// 		a.Debug("Checking for changes")
	// 		a.Changes.SetChanges(a.Git.Changes())
	// 	}
	// }

}
