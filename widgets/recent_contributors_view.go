package widgets

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"strconv"
)

type RecentContributorsView struct {
	name string
	gui  *gocui.Gui
}

func NewRecentContributorsView(gui *gocui.Gui) *RecentContributorsView {
	r := new(RecentContributorsView)
	r.name = "contributors"
	r.gui = gui

	return r
}

func (r *RecentContributorsView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.2 * float64(maxY))
	x2 := int(0.35*float64(maxX)) - 1
	y2 := int(0.4*float64(maxY)) - 1

	v, err := g.SetView(r.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Recent Contributors"

	return nil
}

func (r *RecentContributorsView) UpdateContributors(contributors []*git.Contributor) {
	r.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(r.name)
		if err != nil {
			return err
		}
		v.Clear()

		maxX, _ := v.Size()
		table := view.NewTable(maxX)
		table.AddColumn("NAME", 0.55, view.LEFT)
		table.AddColumn("COMMITS", 0.2, view.RIGHT)
		table.AddColumn("LAST COMMIT", 0.25, view.LEFT)

		for _, contributor := range contributors {
			table.MustAddRow([]string{contributor.Name, strconv.Itoa(contributor.Count), humanize.Time(contributor.LastCommit)})
		}

		table.Render(v)

		return nil
	})
}
