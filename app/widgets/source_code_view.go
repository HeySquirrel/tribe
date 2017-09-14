package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame"
	"github.com/jroimartin/gocui"
)

type SourceCodeView struct {
	name  string
	blame *blame.Blame
	gui   *gocui.Gui
}

func NewSourceCodeView(gui *gocui.Gui, blame *blame.Blame) *SourceCodeView {
	s := new(SourceCodeView)
	s.name = "source"
	s.gui = gui
	s.blame = blame

	return s
}

func (s *SourceCodeView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(0.5*float64(maxX)) - 1
	y2 := int(0.5*float64(maxY)) - 1

	v, err := g.SetView(s.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Source"
	v.Highlight = true
	v.SelBgColor = gocui.ColorCyan
	v.SelFgColor = gocui.ColorBlack

	_, err = g.SetCurrentView(s.name)
	if err != nil {
		return err
	}

	for _, line := range s.blame.Lines {
		fmt.Fprintf(v, "%4d| %s\n", line.Number, line.Text)
	}

	return nil
}
