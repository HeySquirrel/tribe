package widgets

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
	"path/filepath"
)

type SourceCodeView struct {
	name        string
	view        *gocui.View
	presenter   SourcePresenter
	currentLine *model.Line
}

func NewSourceCodeView(presenter SourcePresenter) *SourceCodeView {
	s := new(SourceCodeView)
	s.name = "source"
	s.presenter = presenter

	return s
}

func (s *SourceCodeView) SetCurrentLine(currentLine *model.Line) {
	if s.currentLine == nil {
		s.currentLine = currentLine
		s.view.SetOrigin(0, s.currentLine.Number-1)
	} else {
		moveDistance := currentLine.Number - s.currentLine.Number
		s.currentLine = currentLine
		s.view.MoveCursor(0, moveDistance, false)
	}
}

func (s *SourceCodeView) SetFile(file *model.File) {
	_, title := filepath.Split(file.Name)
	s.view.Title = fmt.Sprintf(" %s:%d,%d ", title, file.Start, file.End)

	for _, line := range file.Lines {
		fmt.Fprintf(s.view, "%3d| %s\n", line.Number, line.Text)
	}
}

func (s *SourceCodeView) Beep() {
	fmt.Print("\a")
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

	s.view = v

	v.Title = "Source"
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlack
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold

	_, err = g.SetCurrentView(s.name)
	if err != nil {
		return err
	}

	return s.setKeyBindings(g)
}

func (s *SourceCodeView) setKeyBindings(g *gocui.Gui) error {
	next := func(g *gocui.Gui, v *gocui.View) error { s.presenter.Next(); return nil }
	previous := func(g *gocui.Gui, v *gocui.View) error { s.presenter.Previous(); return nil }

	err := g.SetKeybinding(s.name, gocui.KeyArrowDown, gocui.ModNone, next)
	if err != nil {
		return err
	}

	err = g.SetKeybinding(s.name, gocui.KeyArrowUp, gocui.ModNone, previous)
	if err != nil {
		return err
	}

	return nil
}
