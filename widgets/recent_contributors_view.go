package widgets

import (
	"github.com/jroimartin/gocui"
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
