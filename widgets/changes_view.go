package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"log"
	"reflect"
)

type SelectionListener func(selectedFile *git.File)

type ChangesView struct {
	name          string
	gui           *gocui.Gui
	listeners     []SelectionListener
	changes       []*git.File
	selectedIndex int
}

func NewChangesView(gui *gocui.Gui) *ChangesView {
	c := new(ChangesView)
	c.name = "changes"
	c.gui = gui
	c.listeners = make([]SelectionListener, 0)
	c.changes = make([]*git.File, 0)
	c.selectedIndex = 0

	return c
}

func (c *ChangesView) AddListener(listener SelectionListener) {
	c.listeners = append(c.listeners, listener)
}

func (c *ChangesView) SetChanges(changes []*git.File) {
	if reflect.DeepEqual(c.changes, changes) {
		return
	}

	c.changes = changes
	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(c.name)
		if err != nil {
			return err
		}
		v.Clear()

		for _, change := range changes {
			fmt.Fprintf(v, "%-30s >%5d >%3d >%3d\n", view.RenderFilename(30, change.Name), len(change.Logs), change.NumberOfDefects(), change.NumberOfStories())
		}

		c.SetSelected(0)
		return nil
	})
}

func (c *ChangesView) GetSelected() *git.File {
	return c.changes[c.selectedIndex]
}

func (c *ChangesView) SetSelected(index int) {
	c.selectedIndex = index

	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(c.name)
		if err != nil {
			return err
		}

		err = v.SetCursor(0, index)
		if err != nil {
			log.Panic(err)
		}

		c.notifyListeners()
		return nil
	})
}

func (c *ChangesView) Next() {
	if c.selectedIndex < len(c.changes)-1 {
		c.SetSelected(c.selectedIndex + 1)
	} else {
		c.SetSelected(0)
	}
}

func (c *ChangesView) Previous() {
	if c.selectedIndex > 0 {
		c.SetSelected(c.selectedIndex - 1)
	} else {
		c.SetSelected(len(c.changes) - 1)
	}
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

	return c.setKeyBindings()
}

func (c *ChangesView) notifyListeners() {
	selected := c.GetSelected()

	for _, listener := range c.listeners {
		listener(selected)
	}
}

func (c *ChangesView) setKeyBindings() error {
	next := func(g *gocui.Gui, v *gocui.View) error { c.Next(); return nil }
	previous := func(g *gocui.Gui, v *gocui.View) error { c.Previous(); return nil }

	err := c.gui.SetKeybinding(c.name, gocui.KeyArrowDown, gocui.ModNone, next)
	if err != nil {
		return err
	}

	err = c.gui.SetKeybinding(c.name, gocui.KeyArrowUp, gocui.ModNone, previous)
	if err != nil {
		return err
	}

	return nil
}
