package app

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	changesView         = "changes"
	contributorsView    = "contributors"
	associatedFilesView = "associatedfiles"
	logsView            = "logs"
	feedView            = "feed"
	relevantWorkView    = "relevantwork"
	legendView          = "legend"
)

type View struct {
	title      string
	text       string
	x1         float64
	y1         float64
	x2         float64
	y2         float64
	highlight  bool
	selBgColor gocui.Attribute
	selFgColor gocui.Attribute
}

var views = map[string]View{
	changesView: {
		title:      "Last Changed",
		text:       "",
		x1:         0.0,
		y1:         0.0,
		x2:         0.3,
		y2:         0.2,
		highlight:  true,
		selBgColor: gocui.ColorCyan,
		selFgColor: gocui.ColorBlack,
	},
	contributorsView: {
		title:     "Frequent Contributors",
		text:      "",
		x1:        0.3,
		y1:        0.0,
		x2:        1.0,
		y2:        0.2,
		highlight: false,
	},
	associatedFilesView: {
		title:     "Associated Files",
		text:      "",
		x1:        0.0,
		y1:        0.2,
		x2:        0.3,
		y2:        0.4,
		highlight: false,
	},
	logsView: {
		title:     "Logs",
		text:      "",
		x1:        0.3,
		y1:        0.2,
		x2:        1.0,
		y2:        0.45,
		highlight: false,
	},
	feedView: {
		title:     "Feed",
		text:      "",
		x1:        0.0,
		y1:        0.4,
		x2:        0.3,
		y2:        0.99,
		highlight: false,
	},
	relevantWorkView: {
		title:     "Relevant Work Items",
		text:      "",
		x1:        0.3,
		y1:        0.45,
		x2:        1.0,
		y2:        0.9,
		highlight: false,
	},
	legendView: {
		title:     "Legend",
		text:      "",
		x1:        0.3,
		y1:        0.9,
		x2:        1.0,
		y2:        0.99,
		highlight: false,
	},
}

var defaultViews = []string{
	changesView,
	contributorsView,
	associatedFilesView,
	logsView,
	feedView,
	relevantWorkView,
	legendView,
}

func (a *App) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	for _, name := range defaultViews {
		view := views[name]

		x1 := int(view.x1 * float64(maxX))
		y1 := int(view.y1 * float64(maxY))
		x2 := int(view.x2*float64(maxX)) - 1
		y2 := int(view.y2*float64(maxY)) - 1

		v, err := g.SetView(name, x1, y1, x2, y2)
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = view.title
		if view.highlight {
			v.Highlight = view.highlight
			v.SelBgColor = view.selBgColor
			v.SelFgColor = view.selFgColor
		}
		fmt.Fprint(v, view.text)
	}

	_, err := g.SetCurrentView(changesView)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) currentFileSelection() string {
	v, err := a.Gui.View(changesView)
	if err != nil {
		log.Panicln(err)
	}

	_, cy := v.Cursor()
	file, err := v.Line(cy)
	if err != nil {
		file = ""
	}

	return file
}

func (a *App) updateLogs(logs string) {
	a.updateView(logsView, func(v *gocui.View) {
		fmt.Fprint(v, logs)
	})
}

func (a *App) UpdateChanges(files []string) {
	a.updateView(changesView, func(v *gocui.View) {
		for _, file := range files {
			fmt.Fprintln(v, file)
		}
		a.updateCurrentFile()
	})
}

func (a *App) updateView(view string, fn func(*gocui.View)) {
	a.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}
		v.Clear()

		fn(v)

		return nil
	})
}
