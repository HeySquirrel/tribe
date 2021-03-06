package widgets

import (
	"fmt"
	"log"

	"github.com/HeySquirrel/tribe/blame/model"
	"github.com/HeySquirrel/tribe/work"
	"github.com/jroimartin/gocui"
	"github.com/kennygrant/sanitize"
)

func NewSourceCodeList(ui *UI) (chan<- *model.File, <-chan *model.Line, gocui.Manager) {
	files := make(chan *model.File)
	selected := make(chan *model.Line)

	l, selections := NewList(ui)

	go func(l *list) {
		for file := range files {
			l.Title(fmt.Sprintf("Source: %s ", file.Name))
			l.SetItems(FileItems(*file), file.Start-1)
			l.Focus()
		}
	}(l)

	go func() {
		for selection := range selections {
			file := model.File(selection.Items.(FileItems))
			selected <- file.GetLine(selection.Index + 1)
		}
	}()

	return files, selected, l
}

func NewItemsList(ui *UI) (chan<- model.Annotation, <-chan *work.FetchedItem, gocui.Manager) {
	annotations := make(chan model.Annotation)
	selected := make(chan *work.FetchedItem)

	l, selections := NewList(ui)

	go func(l *list) {
		for annotation := range annotations {
			workitems := annotation.GetWorkItems()
			l.Title(fmt.Sprintf("Associated Work: %s ", annotation.GetTitle()))
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

			l.Title(fmt.Sprintf("Contributors: %s ", annotation.GetTitle()))
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

			l.Title(fmt.Sprintf("Commits: %s ", annotation.GetTitle()))
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

func NewItemDetails(ui *UI) (chan<- *work.FetchedItem, gocui.Manager) {
	ui.Gui.SetViewOnBottom(ui.Name)
	workitems := make(chan *work.FetchedItem)

	go func() {
		for workitem := range workitems {
			ui.Update(func(v *gocui.View) {
				v.Clear()
				v.Wrap = true

				ui.Title(fmt.Sprintf("%s - ESC to hide", workitem.GetId()))

				fmt.Fprintf(v, "%s - %s\n\n", workitem.GetId(), workitem.GetSummary())
				fmt.Fprintln(v, sanitize.HTML(workitem.GetDescription()))
			})
			hide := ui.Show()
			ui.AddOneUseGlobalKey(gocui.KeyEsc, hide)
		}
	}()

	return workitems, ui
}
