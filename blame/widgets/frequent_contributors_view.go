package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

func NewFileContributorsView(g *gocui.Gui, contributors <-chan *model.AssociatedContributors) gocui.ManagerFunc {
	name := "filecontributors"

	// Handle Updates
	go func(name string) {
		for associated := range contributors {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View(name)
				v.Title = fmt.Sprintf(" Contributors: %s ", associated.Context.GetTitle())

				for _, contributor := range associated.Contributors {
					fmt.Fprintf(v, "  %-20s - %d Commits - %s\n",
						contributor.Name,
						contributor.Count,
						humanize.Time(contributor.LastCommit.Date),
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
		y1 := int(0.75 * float64(maxY))
		x2 := int(0.5*float64(maxX)) - 1
		y2 := int(1.0*float64(maxY)) - 1

		v, err := g.SetView(name, x1, y1, x2, y2)
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Contributors"

		return nil
	}

	return layout
}
