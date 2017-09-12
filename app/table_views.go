package app

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"strconv"
)

func (a *App) UpdateContributors(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		maxX, _ := v.Size()
		table := view.NewTable(maxX)
		table.AddColumn("NAME", 0.55, view.LEFT)
		table.AddColumn("COMMITS", 0.2, view.RIGHT)
		table.AddColumn("LAST COMMIT", 0.25, view.LEFT)

		for _, contributor := range contributors {
			table.MustAddRow([]string{contributor.Name, strconv.Itoa(contributor.Count), humanize.Time(contributor.LastCommit)})
		}

		table.Render(v)
	})
}

func (a *App) UpdateRelatedFiles(files []*git.RelatedFile) {
	a.updateView(associatedFilesView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := view.NewTable(maxX)
		table.AddColumn("NAME", 0.75, view.LEFT)
		table.AddColumn("COMMITS", 0.1, view.RIGHT)
		table.AddColumn("LAST COMMIT", 0.15, view.LEFT)

		for _, file := range files {
			table.MustAddRow([]string{
				view.RenderFilename(file.Name),
				strconv.Itoa(file.Count),
				humanize.Time(file.LastCommit)})
		}

		table.Render(v)
	})
}
