package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

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
