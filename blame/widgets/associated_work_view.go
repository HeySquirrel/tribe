package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

type WorkItemDisplay struct {
	item apis.WorkItem
}

func (w WorkItemDisplay) String() string {
	return fmt.Sprintf("%10s - %s", w.item.GetId(), w.item.GetName())
}

func NewWorkItemsList(ui *UI, works <-chan *model.AssociatedWork) (<-chan apis.WorkItem, gocui.Manager) {
	onSelection := make(chan fmt.Stringer)
	selected := make(chan apis.WorkItem)

	l := NewList(ui, OnEnter, onSelection)
	l.AddGlobalKey(gocui.KeyF2, l.Focus)

	go func(l *list) {
		for work := range works {
			l.Title(fmt.Sprintf(" Associated Work: %s ", work.Context.GetTitle()))
			workitems := make([]fmt.Stringer, len(work.WorkItems))
			for i, item := range work.WorkItems {
				workitems[i] = WorkItemDisplay{item}
			}
			l.SetItems(workitems)
		}
	}(l)

	go func() {
		for item := range onSelection {
			wid := item.(WorkItemDisplay)
			selected <- wid.item
		}
	}()

	return selected, l
}
