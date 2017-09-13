package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis/rally"
	"github.com/jroimartin/gocui"
)

type RelatedWorkView struct {
	name string
	gui  *gocui.Gui
}

func NewRelatedWorkView(gui *gocui.Gui) *RelatedWorkView {
	r := new(RelatedWorkView)
	r.name = "relevantwork"
	r.gui = gui

	return r
}

func (r *RelatedWorkView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.35 * float64(maxX))
	y1 := int(0.45 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(0.9*float64(maxY)) - 1

	v, err := g.SetView(r.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Relevant Work Items"

	return nil
}

func (r *RelatedWorkView) UpdateRelatedWork(workItems []rally.Artifact) {
	r.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(r.name)
		if err != nil {
			return err
		}
		v.Clear()

		for _, workItem := range workItems {
			fmt.Fprintf(v, "%8s - %s\n", workItem.FormattedID, workItem.Name)
		}

		return nil
	})
}
