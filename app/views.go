package app

import (
	"fmt"
	"github.com/heysquirrel/tribe/git"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"log"
	"strconv"
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
		title: "Recent Contributors",
		text:  "",
		x1:    0.3,
		y1:    0.0,
		x2:    1.0,
		y2:    0.2,
	},
	associatedFilesView: {
		title: "Associated Files",
		text:  "",
		x1:    0.0,
		y1:    0.2,
		x2:    0.3,
		y2:    0.4,
	},
	logsView: {
		title: "Logs",
		text:  "",
		x1:    0.3,
		y1:    0.2,
		x2:    1.0,
		y2:    0.45,
	},
	feedView: {
		title: "Feed",
		text:  "",
		x1:    0.0,
		y1:    0.4,
		x2:    0.3,
		y2:    0.99,
	},
	relevantWorkView: {
		title: "Relevant Work Items",
		text:  "",
		x1:    0.3,
		y1:    0.45,
		x2:    1.0,
		y2:    0.9,
	},
	legendView: {
		title: "Legend",
		text:  "",
		x1:    0.3,
		y1:    0.9,
		x2:    1.0,
		y2:    0.99,
	},
	debugView: {
		title:  "Debug",
		text:   "",
		x1:     0.25,
		y1:     0.25,
		x2:     0.7,
		y2:     0.7,
		hidden: true,
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

	_, err := g.SetCurrentView(changesView)
	if err != nil {
		return err
	}

	return a.setKeyBindings()
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

func (a *App) UpdateChanges(files []string) {
	a.updateView(changesView, func(v *gocui.View) {
		for _, file := range files {
			fmt.Fprintln(v, file)
		}
		a.currentFileChanged()
	})
}

func (a *App) UpdateDebug(entries []*tlog.LogEntry) {
	a.updateView(debugView, func(v *gocui.View) {
		table := tablewriter.NewWriter(v)
		table.SetHeader([]string{"Created At", "Message"})
		table.SetBorder(false)

		for _, entry := range entries {
			createdAt := entry.CreatedAt.UTC().Format(time.UnixDate)
			table.Append([]string{createdAt, entry.Message})
		}

		table.Render()
	})
}

func (a *App) UpdateContributors(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		table := tablewriter.NewWriter(v)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, contributor := range contributors {
			table.Append([]string{contributor.Name, strconv.Itoa(contributor.Count), contributor.RelativeDate})
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
