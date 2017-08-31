package app

import (
	"github.com/heysquirrel/tribe/git"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"strconv"
)

type App struct {
	Gui  *gocui.Gui
	Done chan struct{}
	Log  *tlog.Log
}

func New() *App {
	a := new(App)
	a.Done = make(chan struct{})
	a.Log = tlog.New()

	var err error
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
	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) Close() {
	close(a.Done)
	a.Gui.Close()
}

func (a *App) currentFileChanged() {
	file := a.currentFileSelection()
	a.setRecentContributors(git.RecentContributors(file))
}

func (a *App) setRecentContributors(contributors []*git.Contributor) {
	if len(contributors) == 0 {
		return
	}

	a.updateContributors(func(w io.ReadWriter) {
		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, contributor := range contributors {
			table.Append([]string{contributor.Name, strconv.Itoa(contributor.Count), contributor.RelativeDate})
		}

		table.Render()
	})
}
