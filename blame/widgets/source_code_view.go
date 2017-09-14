package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
	"log"
	"path/filepath"
)

type SelectionListener func(selectedLine *model.Line)

type SourceCodeView struct {
	name        string
	gui         *gocui.Gui
	listeners   []SelectionListener
	blame       *model.Blame
	currentLine int
}

func NewSourceCodeView(gui *gocui.Gui, blame *model.Blame) *SourceCodeView {
	s := new(SourceCodeView)
	s.name = "source"
	s.gui = gui
	s.blame = blame
	s.currentLine = 0

	return s
}

func (s *SourceCodeView) AddListener(listener SelectionListener) {
	s.listeners = append(s.listeners, listener)
}

func (s *SourceCodeView) GetSelected() *model.Line {
	return s.blame.Lines[s.currentLine]
}

func (c *SourceCodeView) SetSelected(index int) {
	moveDistance := index - c.currentLine
	c.currentLine = index

	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(c.name)
		if err != nil {
			return err
		}

		v.MoveCursor(0, moveDistance, false)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
}

func (s *SourceCodeView) Next() {
	if s.currentLine < len(s.blame.Lines)-1 {
		s.SetSelected(s.currentLine + 1)
	} else {
		s.SetSelected(0)
	}
}

func (s *SourceCodeView) Previous() {
	if s.currentLine > 0 {
		s.SetSelected(s.currentLine - 1)
	} else {
		s.SetSelected(len(s.blame.Lines) - 1)
	}
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

	_, title := filepath.Split(s.blame.File)
	v.Title = fmt.Sprintf(" %s:%d,%d ", title, s.blame.Start, s.blame.End)
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

	v.SetOrigin(0, s.blame.Start-1)

	return s.setKeyBindings()
}

func (s *SourceCodeView) setKeyBindings() error {
	next := func(g *gocui.Gui, v *gocui.View) error { s.Next(); return nil }
	previous := func(g *gocui.Gui, v *gocui.View) error { s.Previous(); return nil }

	err := s.gui.SetKeybinding(s.name, gocui.KeyArrowDown, gocui.ModNone, next)
	if err != nil {
		return err
	}

	err = s.gui.SetKeybinding(s.name, gocui.KeyArrowUp, gocui.ModNone, previous)
	if err != nil {
		return err
	}

	return nil
}

func (s *SourceCodeView) notifyListeners() {
	selected := s.GetSelected()

	for _, listener := range s.listeners {
		listener(selected)
	}
}