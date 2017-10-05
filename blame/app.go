package blame

import (
	"fmt"
	"log"
	"os"

	"github.com/HeySquirrel/tribe/blame/model"
	"github.com/HeySquirrel/tribe/blame/widgets"
	"github.com/HeySquirrel/tribe/work"
	"github.com/jroimartin/gocui"
)

const (
	newfile          = "newfile"
	newline          = "newline"
	selectedworkitem = "selectedworkitem"
	lineannotation   = "lineannotation"
	fileannotation   = "fileannotation"
)

// App is the entry point into the tribe blame app. It manages the state of all the views
// that make up the entirity of the blame app.
type App struct {
	gui                     *gocui.Gui
	done                    chan struct{}
	annotate                model.Annotate
	views                   []gocui.Manager
	fileListeners           []chan<- *model.File
	fileAnnotationListeners []chan<- model.Annotation
	lineAnnotationListeners []chan<- model.Annotation
	workItemListeners       []chan<- *work.FetchedItem
}

// AddView adds a new view to blame.App. Every view added to the app will be added to
// to the underlying CUI system. blame.App also keeps track of the currently "focused" view and
// allows use of the TAB key to cycle between focusable views.
func (a *App) AddView(view gocui.Manager) {
	a.views = append(a.views, view)
}

// AddFileListener adds a channel that will get notified every time a new file
// is available for display
func (a *App) AddFileListener(c chan<- *model.File) {
	a.fileListeners = append(a.fileListeners, c)
}

// AddFileAnnotationListener adds a channel that will get notified every time a new file based
// annotation is available for display
func (a *App) AddFileAnnotationListener(c chan<- model.Annotation) {
	a.fileAnnotationListeners = append(a.fileAnnotationListeners, c)
}

// AddLineAnnotationListener adds a channel that will get notified every time a new line based
// annotation is available for display
func (a *App) AddLineAnnotationListener(c chan<- model.Annotation) {
	a.lineAnnotationListeners = append(a.lineAnnotationListeners, c)
}

// AddWorkItemListener adds a channel that will get notified every time a new work item
// is available for display
func (a *App) AddWorkItemListener(c chan<- *work.FetchedItem) {
	a.workItemListeners = append(a.workItemListeners, c)
}

// SetHighlightedLine changes the current hightlighted line in blame app. This will trigger any view changes
// that are driven by line changes
func (a *App) SetHighlightedLine(line *model.Line) {
	go func() {
		annotation := a.annotate.Line(line)
		for _, listener := range a.lineAnnotationListeners {
			listener <- annotation
		}
	}()
}

// SetSelectedWorkItem changes the currently selected work item in the blame app
func (a *App) SetSelectedWorkItem(item *work.FetchedItem) {
	go func() {
		for _, listener := range a.workItemListeners {
			listener <- item
		}
	}()
}

// SetFile sets the current file for the blame app
func (a *App) SetFile(filename string, start, end int) {
	go func() {
		file, err := model.NewFile(filename, start, end)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, listener := range a.fileListeners {
			listener <- file
		}

		annotation := a.annotate.File(file)
		for _, listener := range a.fileAnnotationListeners {
			listener <- annotation
		}
	}()
}

// CycleToNextView will Focus the next view in the list of views that can be focused. A view can be focused
// if it contains a FocusOn key
func (a *App) CycleToNextView() {
	var afterFocused bool
	var nextView widgets.Focusable

	for _, view := range a.views {
		ui := view.(widgets.Focusable)
		if ui.CanFocus() && nextView == nil {
			nextView = ui
		}

		if ui.CanFocus() && afterFocused {
			nextView = ui
			break
		}

		if ui.IsFocused() {
			afterFocused = true
			continue
		}
	}

	nextView.Focus()
}

// NewApp creates an instance of the App struct for the given model.File.
// This will create the CUI for blame.
func NewApp(annotate model.Annotate) *App {
	a := &App{done: make(chan struct{}), annotate: annotate}
	var err error

	a.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.gui.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	a.gui.BgColor = gocui.ColorDefault
	a.gui.Highlight = true

	a.addWorkItemDetailView()
	a.addHelpView()
	a.addSourceCodeView()
	a.addFileWorkItemsView()
	a.addFileContributorsView()
	a.addLineCommitsView()
	a.addLineContributorsView()
	a.addLineWorkItemsView()

	a.gui.SetManager(a.views...)

	a.setKeyBindings()

	return a
}

// Loop starts to CUI loop for the blame app. This function will hang until Close is called on this App.
func (a *App) Loop() {
	err := a.gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Close stops the blame app and returns control back to the console
func (a *App) Close() {
	close(a.done)
	a.gui.Close()
}

func (a *App) addSourceCodeView() {
	sourcecode := &widgets.UI{
		Name:    "Source View",
		Startx:  0.0,
		Starty:  0.0,
		Endx:    0.5,
		Endy:    0.5,
		Gui:     a.gui,
		FocusOn: gocui.KeyF1,
	}

	filelistener, linechanges, sourceview := widgets.NewSourceCodeList(sourcecode)

	a.AddView(sourceview)
	a.AddFileListener(filelistener)

	go func() {
		for line := range linechanges {
			a.SetHighlightedLine(line)
		}
	}()
}

func (a *App) addWorkItemDetailView() {
	workitem := &widgets.UI{
		Name:   "WorkItem Detail View",
		Startx: 0.2,
		Starty: 0.2,
		Endx:   0.8,
		Endy:   0.8,
		Gui:    a.gui,
	}

	workitemlistener, workitemview := widgets.NewItemDetails(workitem)
	a.AddView(workitemview)
	a.AddWorkItemListener(workitemlistener)
}

func (a *App) addHelpView() {
	help := &widgets.UI{
		Name:   "Help View",
		Startx: 0.2,
		Starty: 0.2,
		Endx:   0.8,
		Endy:   0.8,
		Gui:    a.gui,
	}

	a.gui.SetViewOnBottom(help.Name)

	show := func() {
		help.Update(func(v *gocui.View) {
			v.Clear()
			help.Title("Shortcuts - F9 to hide")

			help.PrintHelp(v)
			for _, view := range a.views {
				focusable := view.(widgets.Focusable)
				if focusable.CanFocus() {
					documenter := view.(widgets.HelpDocumenter)
					documenter.PrintHelp(v)
				}
			}

			fmt.Fprintln(v, "Global Shortcuts")
			fmt.Fprintf(v, "%10s - Focus next selectable view\n", "Tab")
			fmt.Fprintf(v, "%10s - Quit blame\n", "q")
			fmt.Fprintf(v, "%10s - Quit blame\n", "Ctrl-C")
		})
		hide := help.Show()
		help.AddOneUseGlobalKey(gocui.KeyF9, hide)
	}

	help.AddGlobalKey('h', "Show help", show)
	help.AddGlobalKey('?', "Show help", show)
	help.AddGlobalKey('H', "Show help", show)

	a.AddView(help)
}

func (a *App) addLineCommitsView() {
	commits := &widgets.UI{
		Name:   "Current Line Commit View",
		Startx: 0.5,
		Starty: 0.0,
		Endx:   1.0,
		Endy:   0.4,
		Gui:    a.gui,
	}

	commitlistener, commitview := widgets.NewCommitList(commits)
	a.AddView(commitview)
	a.AddLineAnnotationListener(commitlistener)
}

func (a *App) addLineWorkItemsView() {
	lineworkitems := &widgets.UI{
		Name:    "Current Line WorkItem View",
		Startx:  0.5,
		Starty:  0.4,
		Endx:    1.0,
		Endy:    0.7,
		Gui:     a.gui,
		FocusOn: gocui.KeyF3,
	}

	workitemlistener, selectedworkitems, lineworkview := widgets.NewItemsList(lineworkitems)
	a.AddView(lineworkview)
	a.AddLineAnnotationListener(workitemlistener)

	go func() {
		for workitem := range selectedworkitems {
			a.SetSelectedWorkItem(workitem)
		}
	}()
}

func (a *App) addLineContributorsView() {
	linecontributors := &widgets.UI{
		Name:   "Current Line Contributor View",
		Startx: 0.5,
		Starty: 0.7,
		Endx:   1.0,
		Endy:   1.0,
		Gui:    a.gui,
	}

	contributorlistener, lineconview := widgets.NewContributorsList(linecontributors)
	a.AddView(lineconview)
	a.AddLineAnnotationListener(contributorlistener)
}

func (a *App) addFileWorkItemsView() {
	fileworkitems := &widgets.UI{
		Name:    "Current File WorkItem View",
		Startx:  0.0,
		Starty:  0.5,
		Endx:    0.5,
		Endy:    0.75,
		Gui:     a.gui,
		FocusOn: gocui.KeyF2,
	}

	workitemslistener, selectedworkitems, workview := widgets.NewItemsList(fileworkitems)
	a.AddView(workview)
	a.AddFileAnnotationListener(workitemslistener)

	go func() {
		for workitem := range selectedworkitems {
			a.SetSelectedWorkItem(workitem)
		}
	}()
}

func (a *App) addFileContributorsView() {
	filecontributors := &widgets.UI{
		Name:   "Current File Contributor View",
		Startx: 0.0,
		Starty: 0.75,
		Endx:   0.5,
		Endy:   1.0,
		Gui:    a.gui,
	}

	contributorlistener, conview := widgets.NewContributorsList(filecontributors)
	a.AddView(conview)
	a.AddFileAnnotationListener(contributorlistener)
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

	err = a.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, widgets.ToBinding(a.CycleToNextView))
	if err != nil {
		log.Panicln(err)
	}

	return nil
}
