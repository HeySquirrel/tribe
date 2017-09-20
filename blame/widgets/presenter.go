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
	SetContext(annotation *model.LineAnnotation)
}

type Presenter struct {
	currentLine *model.Line
	file        *model.File
	sourceView  SourceView
	contextView ContextView
	annotate    model.Annotate
}

func NewPresenter(file *model.File, annotate model.Annotate) *Presenter {
	presenter := new(Presenter)
	presenter.file = file
	presenter.currentLine = file.GetLine(file.Start)
	presenter.annotate = annotate

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
		a := p.annotate.Line(p.currentLine)
		p.contextView.SetContext(a)
	}(p)
}
