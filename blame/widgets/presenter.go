package widgets

import (
	"github.com/heysquirrel/tribe/blame/model"
)

type SourcePresenter interface {
	Next()
	Previous()
}

type SourceView interface {
	SetCurrentLine(line *model.Line)
	SetFile(file *model.File)
	Beep()
}

type ContextView interface {
	SetContext(line *model.Line)
}

type Presenter struct {
	currentLine *model.Line
	file        *model.File
	sourceView  SourceView
	contextView ContextView
}

func NewPresenter(file *model.File) *Presenter {
	presenter := new(Presenter)
	presenter.file = file
	presenter.currentLine = file.GetLine(file.Start)

	return presenter
}

func (p *Presenter) SetSourceView(view SourceView) {
	p.sourceView = view

	view.SetFile(p.file)
	view.SetCurrentLine(p.currentLine)
}

func (p *Presenter) SetSourceContextView(view ContextView) {
	p.contextView = view

	p.updateContext()
}

func (p *Presenter) Next() {
	lineNumber := p.currentLine.Number

	if lineNumber < p.file.Len() {
		p.setCurrentLine(p.file.GetLine(lineNumber + 1))
	} else {
		p.sourceView.Beep()
	}
}

func (p *Presenter) Previous() {
	lineNumber := p.currentLine.Number

	if lineNumber > 1 {
		p.setCurrentLine(p.file.GetLine(lineNumber - 1))
	} else {
		p.sourceView.Beep()
	}
}

func (p *Presenter) setCurrentLine(line *model.Line) {
	p.currentLine = line
	p.sourceView.SetCurrentLine(line)
	p.updateContext()
}

func (p *Presenter) updateContext() {
	go func(p *Presenter) {
		p.contextView.SetContext(p.currentLine)
	}(p)
}
