package widgets

import (
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
)

type ThreadSafeSourceView struct {
	gui  *gocui.Gui
	view SourceView
}

type ThreadSafeContextView struct {
	gui  *gocui.Gui
	view ContextView
}

func update(g *gocui.Gui, f func()) {
	g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}

func NewThreadSafeSourceView(gui *gocui.Gui, view SourceView) *ThreadSafeSourceView {
	return &ThreadSafeSourceView{gui: gui, view: view}
}

func (v *ThreadSafeSourceView) SetCurrentLine(line *model.Line) {
	update(v.gui, func() { v.view.SetCurrentLine(line) })
}

func (v *ThreadSafeSourceView) SetFile(file *model.File) {
	update(v.gui, func() { v.view.SetFile(file) })
}

func (v *ThreadSafeSourceView) Beep() {
	update(v.gui, func() { v.view.Beep() })
}

func NewThreadSafeContextView(gui *gocui.Gui, contextView ContextView) *ThreadSafeContextView {
	return &ThreadSafeContextView{gui: gui, view: contextView}
}

func (v *ThreadSafeContextView) SetContext(line *model.History) {
	update(v.gui, func() { v.view.SetContext(line) })
}
