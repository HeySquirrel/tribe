package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

func NewFileWorkItemsView(g *gocui.Gui, works <-chan *model.AssociatedWork) (<-chan apis.WorkItem, gocui.ManagerFunc) {
	name := "fileworkitems"
	selectedWorkItem := make(chan apis.WorkItem)

	// Handle Updates
	go func(name string) {
		for associatedWork := range works {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View(name)
				v.Title = fmt.Sprintf(" Associated Work: %s ", associatedWork.Context.GetTitle())

				for _, item := range associatedWork.WorkItems {
					fmt.Fprintf(v, "%10s - %s\n",
						item.GetId(),
						item.GetName(),
					)
				}

				return nil
			})
		}
	}(name)

	// Initial Layout
	layout := func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		x1 := int(0.0 * float64(maxX))
		y1 := int(0.5 * float64(maxY))
		x2 := int(0.5*float64(maxX)) - 1
		y2 := int(0.75*float64(maxY)) - 1

		v, err := g.SetView(name, x1, y1, x2, y2)
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Associated Work"

		return nil
	}

	return selectedWorkItem, layout
}
