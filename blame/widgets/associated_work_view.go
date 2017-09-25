package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/git"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"regexp"
)

type WorkItems []apis.WorkItem

func (w WorkItems) Display(writer io.Writer) {
	for _, item := range w {
		fmt.Fprintf(writer, "%10s - %s\n", item.GetId(), item.GetName())
	}
}
func (w WorkItems) Len() int { return len(w) }

type ContributorItems []*git.Contributor

func (c ContributorItems) Display(writer io.Writer) {
	for _, contributor := range c {
		fmt.Fprintf(writer, "  %-20s - %d Commits - %s\n",
			contributor.Name,
			contributor.Count,
			humanize.Time(contributor.LastCommit.Date),
		)
	}
}
func (c ContributorItems) Len() int { return len(c) }

type CommitItems []*git.Commit

func (c CommitItems) Display(writer io.Writer) {
	re := regexp.MustCompile("(S|DE|F|s|de|f)[0-9][0-9]+")
	revert := regexp.MustCompile("(r|R)evert")
	magenta := func(s string) string { return color.MagentaString(s) }
	cyan := func(s string) string { return color.CyanString(s) }

	for _, commit := range c {
		subject := re.ReplaceAllStringFunc(commit.Subject, magenta)
		subject = revert.ReplaceAllStringFunc(subject, cyan)

		fmt.Fprintf(writer, " %10s - %s - %s\n",
			commit.Sha[0:9],
			subject,
			humanize.Time(commit.Date),
		)
	}
}
func (c CommitItems) Len() int { return len(c) }

type FileItems model.File

func (f FileItems) Display(writer io.Writer) {
	for _, line := range f.Lines {
		fmt.Fprintf(writer, "%5d| %s\n", line.Number, line.Text)
	}
}
func (f FileItems) Len() int { file := model.File(f); return file.Len() }

func NewSourceCodeList(ui *UI) (chan<- *model.File, <-chan *model.Line, gocui.Manager) {
	files := make(chan *model.File)
	selected := make(chan *model.Line)

	l, selections := NewList(ui)
	l.AddGlobalKey(gocui.KeyF1, l.Focus)

	go func(l *list) {
		for file := range files {
			l.Title(fmt.Sprintf(" Source: %s ", file.Name))
			l.SetItems(FileItems(*file), file.Start-1)
			l.Focus()
		}
	}(l)

	go func() {
		for selection := range selections {
			file := model.File(selection.Items.(FileItems))
			selected <- file.GetLine(selection.Index)
		}
	}()

	return files, selected, l
}

func NewWorkItemsList(ui *UI) (chan<- model.Annotation, <-chan apis.WorkItem, gocui.Manager) {
	annotations := make(chan model.Annotation)
	selected := make(chan apis.WorkItem)

	l, selections := NewList(ui)
	l.AddGlobalKey(gocui.KeyF2, l.Focus)

	go func(l *list) {
		for annotation := range annotations {
			workitems := annotation.GetWorkItems()
			l.Title(fmt.Sprintf(" Associated Work: %s ", annotation.GetTitle()))
			l.SetItems(WorkItems(workitems), 0)
		}
	}(l)

	go func() {
		for selection := range selections {
			if selection.Type != OnEnter {
				continue
			}
			workitems := selection.Items.(WorkItems)
			selected <- workitems[selection.Index]
		}
	}()

	return annotations, selected, l
}

func NewContributorsList(ui *UI) (chan<- model.Annotation, gocui.Manager) {
	annotations := make(chan model.Annotation)

	l, selections := NewList(ui)

	go func(l *list) {
		for annotation := range annotations {
			contributors := annotation.GetContributors()

			l.Title(fmt.Sprintf(" Contributors: %s ", annotation.GetTitle()))
			l.SetItems(ContributorItems(contributors), 0)
		}
	}(l)

	go func() {
		for selection := range selections {
			_, ok := selection.Items.(ContributorItems)
			if !ok {
				log.Panicln("Unknown selection")
			}
		}
	}()

	return annotations, l
}

func NewCommitList(ui *UI) (chan<- model.Annotation, gocui.Manager) {
	annotations := make(chan model.Annotation)

	l, selections := NewList(ui)

	go func(l *list) {
		for annotation := range annotations {
			commits := annotation.GetCommits()

			l.Title(fmt.Sprintf(" Commits: %s ", annotation.GetTitle()))
			l.SetItems(CommitItems(commits), 0)
		}
	}(l)

	go func() {
		for selection := range selections {
			_, ok := selection.Items.(CommitItems)
			if !ok {
				log.Panicln("Unknown selection")
			}
		}
	}()

	return annotations, l
}
