package widgets

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type LegendView struct {
	name string
	gui  *gocui.Gui
}

func NewLegendView(gui *gocui.Gui) *LegendView {
	l := new(LegendView)
	l.name = "legend"
	l.gui = gui

	return l
}

func (l *LegendView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.35 * float64(maxX))
	y1 := int(0.9 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(0.99*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Legend"

	fmt.Fprintln(v, "File > Commits > Defects > Stories")

	return nil
}
