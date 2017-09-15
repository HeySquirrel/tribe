package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

type LineContextView struct {
	name string
	view *gocui.View
}

func NewLineContextView() *LineContextView {
	l := new(LineContextView)
	l.name = "lineview"

	return l
}

func (l *LineContextView) SetCurrentLine(line *model.Line) {
	l.view.Clear()
	l.view.Title = fmt.Sprintf(" Line %d ", line.Number)
	fmt.Fprintln(l.view, line.Text)
}

func (l *LineContextView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.5 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(1.0*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	l.view = v

	v.Title = "Context"

	return nil
}
