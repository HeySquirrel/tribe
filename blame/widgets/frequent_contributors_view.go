package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
)

type FrequentContributorsView struct {
	name         string
	gui          *gocui.Gui
	contributors git.Contributors
}

func NewFrequentContributorsView(gui *gocui.Gui, contributors git.Contributors) *FrequentContributorsView {
	c := new(FrequentContributorsView)
	c.name = "contributors"
	c.gui = gui
	c.contributors = contributors

	return c
}

func (c *FrequentContributorsView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.5 * float64(maxY))
	x2 := int(0.5*float64(maxX)) - 1
	y2 := int(0.75*float64(maxY)) - 1

	v, err := g.SetView(c.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Frequent Contributors"

	for _, contributor := range c.contributors {
		fmt.Fprintf(v, "  %-20s - %d Commits - %s\n",
			contributor.Name,
			contributor.Count,
			humanize.Time(contributor.LastCommit.Date),
		)
	}

	return nil
}
