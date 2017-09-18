package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
	"strings"
)

type LineContextView struct {
	name string
	view *gocui.View
}

func NewLineContextView() *LineContextView {
	l := new(LineContextView)
	l.name = "lineview"

	return l
}

func (l *LineContextView) SetContext(start, end int, commits git.Commits) {
	maxX, _ := l.view.Size()
	maxView := maxX - 2

	l.view.Clear()
	l.view.Title = fmt.Sprintf(" Lines %d - %d ", start, end)

	fmt.Fprintln(l.view, "\n\n  Commits")
	fmt.Fprintf(l.view, "+%s+\n", strings.Repeat("-", maxView))

	for _, commit := range commits {
		fmt.Fprintf(l.view, " %10s - %s - %s\n",
			commit.Sha[0:9],
			commit.Subject,
			humanize.Time(commit.Date),
		)
	}

	fmt.Fprintln(l.view, "\n\n  Work Items")
	fmt.Fprintf(l.view, "+%s+\n", strings.Repeat("-", maxView))

	for _, item := range commits.RelatedWorkItems() {
		fmt.Fprintf(l.view, "  %s\n",
			item,
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
