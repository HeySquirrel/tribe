package app

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"strconv"
)

func (a *App) UpdateContributors(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := tablewriter.NewWriter(v)
		table.SetColWidth(maxX)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, contributor := range contributors {
			table.Append([]string{
				contributor.Name,
				strconv.Itoa(contributor.Count),
				humanize.Time(contributor.LastCommit),
			})
		}

		table.Render()
	})
}

func (a *App) UpdateRelatedFiles(files []*git.RelatedFile) {
	a.updateView(associatedFilesView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := tablewriter.NewWriter(v)
		table.SetColWidth(maxX)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, file := range files {
			table.Append([]string{
				view.RenderFilename(file.Name),
				strconv.Itoa(file.Count),
				humanize.Time(file.LastCommit),
			})
		}

		table.Render()
	})
}
