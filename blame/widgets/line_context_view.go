package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/jroimartin/gocui"
	"log"
	"regexp"
	"strings"
)

type LineContextView struct {
	name   string
	view   *gocui.View
	server apis.WorkItemServer
}

func NewLineContextView(server apis.WorkItemServer) *LineContextView {
	l := new(LineContextView)
	l.name = "lineview"
	l.server = server

	return l
}

func (l *LineContextView) SetContext(line *model.Line) {
	commits := line.GetCommits()

	maxX, _ := l.view.Size()
	maxView := maxX - 2
	re := regexp.MustCompile("(S|DE|F|s|de|f)[0-9][0-9]+")
	revert := regexp.MustCompile("(r|R)evert")

	l.view.Clear()
	l.view.Title = fmt.Sprintf(" Line %d ", line.Number)

	fmt.Fprintln(l.view, "\n\n  Commits")
	fmt.Fprintf(l.view, "+%s+\n", strings.Repeat("-", maxView))

	for _, commit := range commits {
		subject := re.ReplaceAllStringFunc(commit.Subject, func(workitem string) string { return color.MagentaString(workitem) })
		subject = revert.ReplaceAllStringFunc(subject, func(revert string) string { return color.CyanString(revert) })
		fmt.Fprintf(l.view, " %10s - %s - %s\n",
			commit.Sha[0:9],
			subject,
			humanize.Time(commit.Date),
		)
	}

	fmt.Fprintln(l.view, "\n\n  Work Items")
	fmt.Fprintf(l.view, "+%s+\n", strings.Repeat("-", maxView))
	items, err := apis.GetWorkItems(l.server, commits.RelatedWorkItems()...)
	if err != nil {
		log.Panicln(err)
	}

	for _, item := range items {
		fmt.Fprintf(l.view, "%10s - %s\n",
			item.GetId(),
			item.GetName(),
		)
	}

	fmt.Fprintln(l.view, "\n\n  Contributors")
	fmt.Fprintf(l.view, "+%s+\n", strings.Repeat("-", maxView))

	for _, contributor := range commits.RelatedContributors() {
		fmt.Fprintf(l.view, "  %-20s - %d Commits - %s\n",
			contributor.Name,
			contributor.Count,
			humanize.Time(contributor.LastCommit.Date),
		)
	}
}

func (l *LineContextView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.5 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(1.0*float64(maxY)) - 1

	v, err := g.SetView(l.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	l.view = v

	v.Title = "Context"

	return nil
}
