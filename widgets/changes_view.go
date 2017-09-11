package widgets

import (
	"github.com/jroimartin/gocui"
	"log"
)

type ChangesView struct {
	name string
	gui  *gocui.Gui
}

func NewChangesView(gui *gocui.Gui) *ChangesView {
	view := new(ChangesView)
	view.name = "changes"
	view.gui = gui

	return view
}

func (c *ChangesView) GetSelected() string {
	v, err := c.gui.View(c.name)
	if err != nil {
		log.Panicln(err)
	}

	_, cy := v.Cursor()
	file, err := v.Line(cy)
	if err != nil {
		file = ""
	}

	return file
}

func (c *ChangesView) Layout(g *gocui.Gui) error {
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
