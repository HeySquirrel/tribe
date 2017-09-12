package widgets

import (
	"github.com/jroimartin/gocui"
)

type FeedView struct {
	name string
	gui  *gocui.Gui
}

func NewFeedView(gui *gocui.Gui) *FeedView {
	l := new(FeedView)
	l.name = "feed"
	l.gui = gui

	return l
}

func (l *FeedView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.4 * float64(maxY))
	x2 := int(0.35*float64(maxX)) - 1
	y2 := int(0.99*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Feed"

	return nil
}
