package widgets

import (
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"time"
)

type DebugView struct {
	name string
	gui  *gocui.Gui
}

func NewDebugView(gui *gocui.Gui) *DebugView {
	r := new(DebugView)
	r.name = "debug"
	r.gui = gui

	return r
}

func (d *DebugView) Hide() {
	d.gui.SetViewOnBottom("debug")
}

func (r *DebugView) UpdateDebug(entries []*tlog.LogEntry) {
	r.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(r.name)
		if err != nil {
			return err
		}
		v.Clear()

		maxX, _ := v.Size()
		table := view.NewTable(maxX)
		table.AddColumn("CREATED AT", 0.3, view.LEFT)
		table.AddColumn("MESSAGE", 0.7, view.LEFT)

		for _, entry := range entries {
			createdAt := entry.CreatedAt.UTC().Format(time.UnixDate)
			table.MustAddRow([]string{createdAt, entry.Message})
		}

		table.Render(v)

		return nil
	})
}

func (d *DebugView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.25 * float64(maxX))
	y1 := int(0.25 * float64(maxY))
	x2 := int(0.8*float64(maxX)) - 1
	y2 := int(0.8*float64(maxY)) - 1

	v, err := g.SetView(d.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Debug"

	g.SetViewOnBottom(d.name)

	return d.setKeyBindings()
}

func (d *DebugView) setKeyBindings() error {
	hide := func(g *gocui.Gui, v *gocui.View) error { d.Hide(); return nil }
	err := d.gui.SetKeybinding(d.name, gocui.KeyF1, gocui.ModNone, hide)
	if err != nil {
		return err
	}
	return nil
}
