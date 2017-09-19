package blame

import (
	"github.com/heysquirrel/tribe/apis/rally"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/blame/widgets"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
)

type BlameApp struct {
	Gui       *gocui.Gui
	Done      chan struct{}
	Presenter *widgets.Presenter
}

func NewBlameApp(blame *model.File) *BlameApp {
	usr, err := user.Current()
	if err != nil {
		log.Panicln(err)
	}

	configFile := filepath.Join(usr.HomeDir, ".tribe")
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicln(err)
	}

	api := rally.New(string(config))

	a := new(BlameApp)
	a.Done = make(chan struct{})

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.Presenter = widgets.NewPresenter(blame)
	source := widgets.NewSourceCodeView(a.Presenter)
	lineContext := widgets.NewLineContextView(api)

	a.Presenter.SetSourceView(widgets.NewThreadSafeSourceView(a.Gui, source))
	a.Presenter.SetSourceContextView(widgets.NewThreadSafeContextView(a.Gui, lineContext))

	logs, err := blame.Logs()
	if err != nil {
		log.Panicln(err)
	}

	contributors := logs.RelatedContributors()
	workItems := logs.RelatedWorkItems()

	a.Gui.SetManager(
		source,
		widgets.NewFrequentContributorsView(a.Gui, contributors),
		widgets.NewAssociatedWorkView(a.Gui, workItems),
		lineContext,
	)

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
