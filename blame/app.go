package blame

import (
	"log"

	"github.com/HeySquirrel/tribe/blame/model"
	"github.com/HeySquirrel/tribe/blame/widgets"
	"github.com/HeySquirrel/tribe/work"
	"github.com/jroimartin/gocui"
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

	workitem := &widgets.UI{
		Name:   "workitem",
		Startx: 0.2,
		Starty: 0.2,
		Endx:   0.8,
		Endy:   0.8,
		Gui:    a.Gui,
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
		Name:    "lineworkitems",
		Startx:  0.5,
		Starty:  0.4,
		Endx:    1.0,
		Endy:    0.7,
		Gui:     a.Gui,
		FocusOn: gocui.KeyF3,
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

	workitems := make(chan *work.FetchedItem)

	filein, lineout, sourceview := widgets.NewSourceCodeList(sourcecode)

	commitin, commitview := widgets.NewCommitList(commits)
	lineworkin, lineworkout, lineworkview := widgets.NewItemsList(lineworkitems)
	lineconin, lineconview := widgets.NewContributorsList(linecontributors)

	workin, fileworkout, workview := widgets.NewItemsList(fileworkitems)
	conin, conview := widgets.NewContributorsList(filecontributors)

	workitemview := widgets.NewItemDetails(workitem, workitems)

	a.Gui.SetManager(
		workitemview,
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
		for {
			select {
			case lineout := <-lineworkout:
				workitems <- lineout
			case fileout := <-fileworkout:
				workitems <- fileout
			}
		}
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
