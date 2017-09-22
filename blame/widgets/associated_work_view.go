package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
	"log"
)

type WorkItemDisplay struct {
	item apis.WorkItem
}

func (w WorkItemDisplay) String() string {
	return fmt.Sprintf("%10s - %s", w.item.GetId(), w.item.GetName())
}

type ContributorDisplay struct {
	contributor *git.Contributor
}

func (c ContributorDisplay) String() string {
	return fmt.Sprintf("  %-20s - %d Commits - %s",
		c.contributor.Name,
		c.contributor.Count,
		humanize.Time(c.contributor.LastCommit.Date),
	)
}

func NewWorkItemsList(ui *UI) (chan<- model.Annotation, <-chan apis.WorkItem, gocui.Manager) {
	annotations := make(chan model.Annotation)
	onSelection := make(chan fmt.Stringer)
	selected := make(chan apis.WorkItem)

	l := NewList(ui, OnEnter, onSelection)
	l.AddGlobalKey(gocui.KeyF2, l.Focus)

	go func(l *list) {
		for annotation := range annotations {
			workitems := annotation.GetWorkItems()

			l.Title(fmt.Sprintf(" Associated Work: %s ", annotation.GetTitle()))
			displays := make([]fmt.Stringer, len(workitems))
			for i, item := range workitems {
				displays[i] = WorkItemDisplay{item}
			}
			l.SetItems(displays)
		}
	}(l)

	go func() {
		for item := range onSelection {
			wid := item.(WorkItemDisplay)
			selected <- wid.item
		}
	}()

	return annotations, selected, l
}

func NewContributorsList(ui *UI) (chan<- model.Annotation, gocui.Manager) {
	annotations := make(chan model.Annotation)
	onSelection := make(chan fmt.Stringer)

	l := NewList(ui, OnEnter, onSelection)

	go func(l *list) {
		for annotation := range annotations {
			contributors := annotation.GetContributors()

			l.Title(fmt.Sprintf(" Contributors: %s ", annotation.GetTitle()))
			displays := make([]fmt.Stringer, len(contributors))
			for i, item := range contributors {
				displays[i] = ContributorDisplay{item}
			}
			l.SetItems(displays)
		}
	}(l)

	go func() {
		for item := range onSelection {
			_, ok := item.(ContributorDisplay)
			if !ok {
				log.Panicln("Unknown selection")
			}
		}
	}()

	return annotations, l
}
