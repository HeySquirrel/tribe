package widgets

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"strconv"
)

type AssociatedFilesView struct {
	name string
	gui  *gocui.Gui
}

func NewAssociatedFilesView(gui *gocui.Gui) *AssociatedFilesView {
	a := new(AssociatedFilesView)
	a.name = "associatedfiles"
	a.gui = gui

	return a
}

func (a *AssociatedFilesView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.35 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(1.0*float64(maxX)) - 1
	y2 := int(0.2*float64(maxY)) - 1

	v, err := g.SetView(a.name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	v.Title = "Associated Files"

	return nil
}

func (a *AssociatedFilesView) UpdateRelatedFiles(files []*git.RelatedFile) {
	a.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(a.name)
		if err != nil {
			return err
		}
		v.Clear()

		maxX, _ := v.Size()

		table := view.NewTable(maxX)
		table.AddColumn("NAME", 0.75, view.LEFT)
		table.AddColumn("COMMITS", 0.1, view.RIGHT)
		table.AddColumn("LAST COMMIT", 0.15, view.LEFT)

		for _, file := range files {
			table.MustAddRow([]string{
				view.RenderFilename(file.Name),
				strconv.Itoa(file.Count),
				humanize.Time(file.LastCommit)})
		}

		table.Render(v)

		return nil
	})
}
