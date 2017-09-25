package blame

import (
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/blame/widgets"
	"github.com/jroimartin/gocui"
	"log"
)

type BlameApp struct {
	Gui  *gocui.Gui
	Done chan struct{}
}

func NewBlameApp(file *model.File, annotate model.Annotate) *BlameApp {
	a := new(BlameApp)
	a.Done = make(chan struct{})
	var err error

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.Gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	a.Gui.BgColor = gocui.ColorDefault
	a.Gui.Highlight = true

	sourcecode := &widgets.UI{
		Name:    "source",
		Startx:  0.0,
		Starty:  0.0,
		Endx:    0.5,
		Endy:    0.5,
		Gui:     a.Gui,
		FocusOn: gocui.KeyF1,
	}

	commits := &widgets.UI{
		Name:   "commits",
		Startx: 0.5,
		Starty: 0.0,
		Endx:   1.0,
		Endy:   0.4,
		Gui:    a.Gui,
	}

	lineworkitems := &widgets.UI{
		Name:   "lineworkitems",
		Startx: 0.5,
		Starty: 0.4,
		Endx:   1.0,
		Endy:   0.7,
		Gui:    a.Gui,
	}

	linecontributors := &widgets.UI{
		Name:   "linecontributors",
		Startx: 0.5,
		Starty: 0.7,
		Endx:   1.0,
		Endy:   1.0,
		Gui:    a.Gui,
	}

	fileworkitems := &widgets.UI{
		Name:    "fileworkitems",
		Startx:  0.0,
		Starty:  0.5,
		Endx:    0.5,
		Endy:    0.75,
		Gui:     a.Gui,
		FocusOn: gocui.KeyF2,
	}

	filecontributors := &widgets.UI{
		Name:   "filecontributors",
		Startx: 0.0,
		Starty: 0.75,
		Endx:   0.5,
		Endy:   1.0,
		Gui:    a.Gui,
	}

	filein, lineout, sourceview := widgets.NewSourceCodeList(sourcecode)

	commitin, commitview := widgets.NewCommitList(commits)
	lineworkin, _, lineworkview := widgets.NewWorkItemsList(lineworkitems)
	lineconin, lineconview := widgets.NewContributorsList(linecontributors)

	workin, _, workview := widgets.NewWorkItemsList(fileworkitems)
	conin, conview := widgets.NewContributorsList(filecontributors)

	a.Gui.SetManager(
		sourceview,
		commitview,
		lineworkview,
		lineconview,
		workview,
		conview,
	)

	go func() {
		filein <- file
	}()

	go func() {
		for line := range lineout {
			annotation := annotate.Line(line)
			commitin <- annotation
			lineworkin <- annotation
			lineconin <- annotation
		}
	}()

	go func() {
		annotation := annotate.File(file)
		workin <- annotation
		conin <- annotation
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
