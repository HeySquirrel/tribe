package widgets

import (
	"github.com/jroimartin/gocui"
)

type Changes struct {
	name string
}

func NewChanges() *Changes {
	changes := new(Changes)
	changes.name = "changes"

	return changes
}

func (c *Changes) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(0.35*float64(maxX)) - 1
	y2 := int(0.2*float64(maxY)) - 1

	v, err := g.SetView(c.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Last Changed"
	v.Highlight = true
	v.SelBgColor = gocui.ColorCyan
	v.SelFgColor = gocui.ColorBlack

	_, err = g.SetCurrentView(c.name)
	if err != nil {
		return err
	}

	return nil
}
