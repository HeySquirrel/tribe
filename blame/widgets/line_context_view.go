package widgets

import (
	"github.com/jroimartin/gocui"
)

type LineContextView struct {
	name string
	gui  *gocui.Gui
}

func NewLineContextView(gui *gocui.Gui) *LineContextView {
	l := new(LineContextView)
	l.name = "lineview"
	l.gui = gui

	return l
}

func (l *LineContextView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.5 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(1.0*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Line 104"

	return nil
}
