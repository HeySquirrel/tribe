package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

func ToBinding(f func()) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error { f(); return nil }
}

type KeyBinding struct {
	view    string
	key     interface{}
	handler func()
}

type UI struct {
	name   string
	startx float64
	starty float64
	endx   float64
	endy   float64
	gui    *gocui.Gui
	keys   []KeyBinding
}

func (u *UI) Update(f func(v *gocui.View)) {
	u.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(u.name)
		if err != nil {
			return err
		}

		f(v)
		return nil
	})
}

func (u *UI) Title(title string) {
	u.Update(func(v *gocui.View) {
		v.Title = title
	})
}

func (u *UI) Focus() {
	u.Update(func(v *gocui.View) {
		u.gui.CurrentView().Highlight = false
		v.Highlight = true
		u.gui.SetCurrentView(u.name)
	})
}

func (u *UI) AddLocalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, KeyBinding{u.name, key, handler})
}

func (u *UI) AddGlobalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, KeyBinding{"", key, handler})
}

func (u *UI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(u.startx * float64(maxX))
	y1 := int(u.starty * float64(maxY))
	x2 := int(u.endx*float64(maxX)) - 1
	y2 := int(u.endy*float64(maxY)) - 1

	_, err := g.SetView(u.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	return u.registerKeyBindings(g)
}

func (u *UI) registerKeyBindings(g *gocui.Gui) error {
	for _, binding := range u.keys {
		err := g.SetKeybinding(binding.view, binding.key, gocui.ModNone, ToBinding(binding.handler))
		if err != nil {
			return err
		}
	}
	return nil
}

type list struct {
	*UI
	items    apis.WorkItems
	current  int
	selected chan apis.WorkItem
}

func (l *list) setSelection(index int) {
	if len(l.items) == 0 {
		return
	}

	if index < 0 || index >= len(l.items) {
		fmt.Print("\a")
		return
	}

	l.Update(func(v *gocui.View) {
		if l.current == -1 {
			l.current = 0
			v.SetOrigin(0, 0)
		} else {
			moveDistance := index - l.current
			l.current = index
			v.MoveCursor(0, moveDistance, false)
		}
	})
}

func (l *list) setItems(items apis.WorkItems) {
	l.Update(func(v *gocui.View) {
		v.Clear()

		for _, item := range items {
			fmt.Fprintf(v, "%10s - %s\n",
				item.GetId(),
				item.GetName(),
			)
		}

		l.items = items
		l.setSelection(0)
	})
}

func (l *list) next()     { l.setSelection(l.current + 1) }
func (l *list) previous() { l.setSelection(l.current - 1) }

func NewFileWorkItemsView(g *gocui.Gui, works <-chan *model.AssociatedWork) (<-chan apis.WorkItem, gocui.Manager) {
	selected := make(chan apis.WorkItem)
	ui := &UI{
		name:   "fileworkitems",
		startx: 0.0,
		starty: 0.5,
		endx:   0.5,
		endy:   0.75,
		gui:    g,
		keys:   make([]KeyBinding, 0),
	}
	l := &list{ui, []apis.WorkItem{}, -1, selected}
	l.AddLocalKey(gocui.KeyArrowUp, l.previous)
	l.AddLocalKey(gocui.KeyArrowDown, l.next)
	l.AddGlobalKey(gocui.KeyF2, l.Focus)

	go func(l *list) {
		for work := range works {
			l.Title(fmt.Sprintf(" Associated Work: %s ", work.Context.GetTitle()))
			l.setItems(work.WorkItems)
		}
	}(l)

	return selected, l
}
