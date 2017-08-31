package app

import (
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"strconv"
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

func (a *App) currentFileChanged() {
	file := a.currentFileSelection()
	a.setFrequentContributors(git.FrequentContributors(file))
}

func (a *App) setFrequentContributors(contributors []*git.Contributor) {
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
