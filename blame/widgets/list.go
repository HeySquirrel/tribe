package widgets

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type SelectionEvent int

type list struct {
	*UI
	items      []fmt.Stringer
	current    int
	selectFire SelectionEvent
	selected   chan fmt.Stringer
}

const (
	OnSelect SelectionEvent = iota
	OnEnter
)

func NewList(ui *UI, selectFire SelectionEvent, selected chan fmt.Stringer) *list {
	l := &list{
		ui,
		make([]fmt.Stringer, 0),
		-1,
		selectFire,
		selected,
	}

	l.AddLocalKey(gocui.KeyArrowUp, l.Previous)
	l.AddLocalKey(gocui.KeyArrowDown, l.Next)
	l.AddLocalKey(gocui.KeyEnter, func() { l.fire(OnEnter) })

	return l
}

func (l *list) SetSelection(index int) {
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
		l.fire(OnSelect)
	})
}

func (l *list) fire(event SelectionEvent) {
	if event != l.selectFire {
		return
	}

	go func() {
		l.selected <- l.items[l.current]
	}()
}

func (l *list) SetItems(items []fmt.Stringer) {
	l.Update(func(v *gocui.View) {
		v.Clear()

		for _, item := range items {
			fmt.Fprintln(v, item)
		}

		l.items = items
		l.SetSelection(0)
	})
}

func (l *list) Next()     { l.SetSelection(l.current + 1) }
func (l *list) Previous() { l.SetSelection(l.current - 1) }
