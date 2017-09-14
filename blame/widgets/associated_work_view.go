package widgets

import (
	"github.com/jroimartin/gocui"
)

type AssociatedWorkView struct {
	name string
	gui  *gocui.Gui
}

func NewAssociatedWorkView(gui *gocui.Gui) *AssociatedWorkView {
	a := new(AssociatedWorkView)
	a.name = "workitems"
	a.gui = gui

	return a
}

func (a *AssociatedWorkView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.75 * float64(maxY))
	x2 := int(0.5*float64(maxX)) - 1
	y2 := int(1.0*float64(maxY)) - 1

	v, err := g.SetView(a.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Associated Work"

	return nil
}
