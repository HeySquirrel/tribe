package widgets

import (
	"github.com/jroimartin/gocui"
)

type LogsView struct {
	name string
	gui  *gocui.Gui
}

func NewLogsView(gui *gocui.Gui) *LogsView {
	l := new(LogsView)
	l.name = "logs"
	l.gui = gui

	return l
}

func (l *LogsView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.35 * float64(maxX))
	y1 := int(0.2 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(0.45*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Logs"

	return nil
}
