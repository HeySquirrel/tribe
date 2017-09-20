package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/jroimartin/gocui"
)

type AssociatedWorkView struct {
	name      string
	gui       *gocui.Gui
	workItems []apis.WorkItem
}

func NewAssociatedWorkView(gui *gocui.Gui, workItems []apis.WorkItem) *AssociatedWorkView {
	a := new(AssociatedWorkView)
	a.name = "workitems"
	a.gui = gui
	a.workItems = workItems

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

	for _, item := range a.workItems {
		fmt.Fprintf(v, "%10s - %s\n",
			item.GetId(),
			item.GetName(),
		)
	}

	return nil
}
