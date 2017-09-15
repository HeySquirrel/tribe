package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

type LineContextView struct {
	name        string
	gui         *gocui.Gui
	currentLine *model.Line
}

func NewLineContextView(gui *gocui.Gui, currentLine *model.Line) *LineContextView {
	l := new(LineContextView)
	l.name = "lineview"
	l.gui = gui
	l.currentLine = currentLine

	return l
}

func (l *LineContextView) SetCurrentLine(currentLine *model.Line) {
	l.currentLine = currentLine

	l.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(l.name)
		if err != nil {
			return err
		}
		v.Clear()

		v.Title = fmt.Sprintf(" Line %d ", l.currentLine.Number)
		fmt.Fprintln(v, l.currentLine.Text)

		return nil
	})
}

func (l *LineContextView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.5 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(1.0*float64(maxY)) - 1

	_, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	l.SetCurrentLine(l.currentLine)

	return nil
}
