package blame

import (
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/blame/widgets"
	"github.com/jroimartin/gocui"
	"log"
)

type BlameApp struct {
	Gui       *gocui.Gui
	Done      chan struct{}
	Presenter *widgets.Presenter
}

func NewBlameApp(file *model.File, annotate model.Annotate) *BlameApp {
	a := new(BlameApp)
	a.Done = make(chan struct{})
	var err error

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	associatedWork := make(chan *model.AssociatedWork)
	associatedContributors := make(chan *model.AssociatedContributors)

	a.Gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	a.Gui.BgColor = gocui.ColorDefault
	a.Gui.Highlight = true

	a.Presenter = widgets.NewPresenter(file, annotate)
	source := widgets.NewSourceCodeView(a.Presenter)
	lineContext := widgets.NewLineContextView()

	a.Presenter.SetSourceView(widgets.NewThreadSafeSourceView(a.Gui, source))
	a.Presenter.SetSourceContextView(widgets.NewThreadSafeContextView(a.Gui, lineContext))

	ui := &widgets.UI{
		Name:   "fileworkitems",
		Startx: 0.0,
		Starty: 0.5,
		Endx:   0.5,
		Endy:   0.75,
		Gui:    a.Gui,
	}

	_, workview := widgets.NewWorkItemsList(ui, associatedWork)

	a.Gui.SetManager(
		source,
		workview,
		widgets.NewFileContributorsView(a.Gui, associatedContributors),
		lineContext,
	)

	go func() {
		annotation := annotate.File(file)
		associatedWork <- &model.AssociatedWork{annotation, annotation.GetWorkItems()}
		associatedContributors <- &model.AssociatedContributors{annotation, annotation.GetContributors()}
	}()

	a.setKeyBindings()

	return a
}

func (a *BlameApp) Loop() {
	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *BlameApp) Close() {
	close(a.Done)
	a.Gui.Close()
}

func (a *BlameApp) setKeyBindings() error {
	quit := func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }
	err := a.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}
