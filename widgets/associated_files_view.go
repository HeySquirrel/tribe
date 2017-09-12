package widgets

import (
	"github.com/jroimartin/gocui"
)

type AssociatedFilesView struct {
	name string
	gui  *gocui.Gui
}

func NewAssociatedFilesView(gui *gocui.Gui) *AssociatedFilesView {
	a := new(AssociatedFilesView)
	a.name = "associatedfiles"
	a.gui = gui

	return a
}

func (a *AssociatedFilesView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.35 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(0.2*float64(maxY)) - 1

	v, err := g.SetView(a.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Associated Files"

	return nil
}
