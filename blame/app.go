package blame

import (
	"log"

	"github.com/HeySquirrel/tribe/blame/model"
	"github.com/HeySquirrel/tribe/blame/widgets"
	"github.com/cskr/pubsub"
	"github.com/jroimartin/gocui"
)

// App is the entry point into the tribe blame app
type App struct {
	gui    *gocui.Gui
	done   chan struct{}
	pubsub *pubsub.PubSub
}

// NewApp creates an instance of the App struct for the given model.File.
// This will create the CUI for blame.
func NewApp(file *model.File, annotate model.Annotate) *App {
	a := new(App)
	a.done = make(chan struct{})
	a.pubsub = pubsub.New(1)
	var err error

	a.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	a.gui.BgColor = gocui.ColorDefault
	a.gui.Highlight = true

	sourcecode := &widgets.UI{
		Name:    "source",
		Startx:  0.0,
		Starty:  0.0,
		Endx:    0.5,
		Endy:    0.5,
		Gui:     a.gui,
		FocusOn: gocui.KeyF1,
	}

	workitem := &widgets.UI{
		Name:   "workitem",
		Startx: 0.2,
		Starty: 0.2,
		Endx:   0.8,
		Endy:   0.8,
		Gui:    a.gui,
	}

	commits := &widgets.UI{
		Name:   "commits",
		Startx: 0.5,
		Starty: 0.0,
		Endx:   1.0,
		Endy:   0.4,
		Gui:    a.gui,
	}

	lineworkitems := &widgets.UI{
		Name:    "lineworkitems",
		Startx:  0.5,
		Starty:  0.4,
		Endx:    1.0,
		Endy:    0.7,
		Gui:     a.gui,
		FocusOn: gocui.KeyF3,
	}

	linecontributors := &widgets.UI{
		Name:   "linecontributors",
		Startx: 0.5,
		Starty: 0.7,
		Endx:   1.0,
		Endy:   1.0,
		Gui:    a.gui,
	}

	fileworkitems := &widgets.UI{
		Name:    "fileworkitems",
		Startx:  0.0,
		Starty:  0.5,
		Endx:    0.5,
		Endy:    0.75,
		Gui:     a.gui,
		FocusOn: gocui.KeyF2,
	}

	filecontributors := &widgets.UI{
		Name:   "filecontributors",
		Startx: 0.0,
		Starty: 0.75,
		Endx:   0.5,
		Endy:   1.0,
		Gui:    a.gui,
	}

	filein, lineout, sourceview := widgets.NewSourceCodeList(sourcecode)

	commitin, commitview := widgets.NewCommitList(commits)
	lineworkin, lineworkout, lineworkview := widgets.NewItemsList(lineworkitems)
	lineconin, lineconview := widgets.NewContributorsList(linecontributors)

	workin, fileworkout, workview := widgets.NewItemsList(fileworkitems)
	conin, conview := widgets.NewContributorsList(filecontributors)

	workitems, workitemview := widgets.NewItemDetails(workitem)

	a.gui.SetManager(
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

func (a *App) Loop() {
	err := a.gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) Close() {
	close(a.done)
	a.gui.Close()
}

func (a *App) setKeyBindings() error {
	quit := func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }
	err := a.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}

	err = a.gui.SetKeybinding("", 'q', gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}
