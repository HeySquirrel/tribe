package app

import (
	"fmt"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"time"
)

const (
	changesView         = "changes"
	contributorsView    = "contributors"
	associatedFilesView = "associatedfiles"
	logsView            = "logs"
	feedView            = "feed"
	relevantWorkView    = "relevantwork"
	legendView          = "legend"
	debugView           = "debug"
)

type View struct {
	title      string
	text       string
	x1         float64
	y1         float64
	x2         float64
	y2         float64
	highlight  bool
	hidden     bool
	selBgColor gocui.Attribute
	selFgColor gocui.Attribute
}

var views = map[string]View{
	feedView: {
		title: "Feed",
		text:  "",
		x1:    0.0,
		y1:    0.4,
		x2:    0.35,
		y2:    0.99,
	},
	legendView: {
		title: "Legend",
		text:  "",
		x1:    0.35,
		y1:    0.9,
		x2:    1.0,
		y2:    0.99,
	},
	debugView: {
		title:  "Debug",
		text:   "",
		x1:     0.25,
		y1:     0.25,
		x2:     0.8,
		y2:     0.8,
		hidden: true,
	},
}

var defaultViews = []string{
	feedView,
	legendView,
	debugView,
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

		if view.hidden {
			g.SetViewOnBottom(name)
		}
	}

	return a.setKeyBindings()
}

func (a *App) UpdateDebug(entries []*tlog.LogEntry) {
	a.updateView(debugView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := tablewriter.NewWriter(v)
		table.SetColWidth(maxX)
		table.SetHeader([]string{"Created At", "Message"})
		table.SetBorder(false)

		for _, entry := range entries {
			createdAt := entry.CreatedAt.UTC().Format(time.UnixDate)
			table.Append([]string{createdAt, entry.Message})
		}

		table.Render()
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
