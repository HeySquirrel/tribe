package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
	"log"
	"regexp"
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

type CommitDisplay struct {
	commit *git.Commit
}

func (c CommitDisplay) String() string {
	commit := c.commit
	re := regexp.MustCompile("(S|DE|F|s|de|f)[0-9][0-9]+")
	revert := regexp.MustCompile("(r|R)evert")

	subject := re.ReplaceAllStringFunc(commit.Subject, func(workitem string) string { return color.MagentaString(workitem) })
	subject = revert.ReplaceAllStringFunc(subject, func(revert string) string { return color.CyanString(revert) })
	return fmt.Sprintf(" %10s - %s - %s",
		commit.Sha[0:9],
		subject,
		humanize.Time(commit.Date),
	)
}

func NewSourceCodeList(ui *UI) (chan<- *model.File, <-chan *model.Line, gocui.Manager) {
	files := make(chan *model.File)
	onSelection := make(chan fmt.Stringer)
	selected := make(chan *model.Line)

	l := NewList(ui, OnSelect, onSelection)
	l.AddGlobalKey(gocui.KeyF1, l.Focus)

	go func(l *list) {
		for file := range files {
			lines := file.Lines

			l.Title(fmt.Sprintf(" Source: %s ", file.Name))
			displays := make([]fmt.Stringer, len(lines))
			for i, line := range lines {
				displays[i] = line
			}
			l.SetItems(displays)
			l.SetSelection(file.Start - 1)
			l.Focus()
		}
	}(l)

	go func() {
		for line := range onSelection {
			l := line.(*model.Line)
			selected <- l
		}
	}()

	return files, selected, l
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

func NewCommitList(ui *UI) (chan<- model.Annotation, gocui.Manager) {
	annotations := make(chan model.Annotation)
	onSelection := make(chan fmt.Stringer)

	l := NewList(ui, OnEnter, onSelection)

	go func(l *list) {
		for annotation := range annotations {
			commits := annotation.GetCommits()

			l.Title(fmt.Sprintf(" Commits: %s ", annotation.GetTitle()))
			displays := make([]fmt.Stringer, len(commits))
			for i, item := range commits {
				displays[i] = CommitDisplay{item}
			}
			l.SetItems(displays)
		}
	}(l)

	go func() {
		for item := range onSelection {
			_, ok := item.(CommitDisplay)
			if !ok {
				log.Panicln("Unknown selection")
			}
		}
	}()

	return annotations, l
}
