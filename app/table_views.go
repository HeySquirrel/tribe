package app

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"io"
	"strconv"
	"strings"
)

type Column struct {
	name string
	size int
}

type Row []string

type Table struct {
	width   int
	columns []Column
	rows    []Row
}

func NewColumn(name string, size int) Column {
	return Column{name: name, size: size}
}

func NewTable(width int) *Table {
	table := new(Table)
	table.width = width
	table.columns = make([]Column, 0)
	table.rows = make([]Row, 0)

	return table
}

func (t *Table) AddColumn(name string, size float64) {
	columnSize := int(float64(t.width) * size)
	t.columns = append(t.columns, NewColumn(name, columnSize))
}

func (t *Table) MustAddRow(row Row) {
	if len(t.columns) != len(row) {
		panic("Row size should match column size")
	}

	t.rows = append(t.rows, row)
}

func Center2(s string, width int) string {
	leftPad := width/2 + len(s)/2

	if leftPad%2 != 0 {
		leftPad = leftPad + 1
	}

	return fmt.Sprintf(fmt.Sprintf("%%-%ds", width), fmt.Sprintf(fmt.Sprintf("%%%ds", leftPad), s))
}

func (t *Table) Render(w io.Writer) {
	maxView := t.width - 2

	header := make([]string, 0)
	for _, column := range t.columns {
		header = append(header, Center2(column.name, column.size))
	}
	fmt.Fprintln(w, strings.Join(header, "|"))

	fmt.Fprintf(w, "+%s+\n", strings.Repeat("-", maxView))

	for i, row := range t.rows {
		columns := make([]string, 0)

		for j, column := range t.columns {
			data := row[j]
			if i == 0 && j == 0 {
				data = fmt.Sprintf(" 🌶  %s", data)
			}
			columnFormat := fmt.Sprintf(" %%-%ds", column.size-1)
			columns = append(columns, fmt.Sprintf(columnFormat, data))
		}

		fmt.Fprintln(w, strings.Join(columns, "|"))
	}
}

func (a *App) UpdateContributors2(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		maxX, _ := v.Size()
		table := NewTable(maxX)
		table.AddColumn("NAME", 0.55)
		table.AddColumn("COMMITS", 0.2)
		table.AddColumn("LAST COMMIT", 0.25)

		for _, contributor := range contributors {
			table.MustAddRow([]string{contributor.Name, strconv.Itoa(contributor.Count), humanize.Time(contributor.LastCommit)})
		}

		table.Render(v)
	})
}

func (a *App) UpdateRelatedFiles(files []*git.RelatedFile) {
	a.updateView(associatedFilesView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := NewTable(maxX)
		table.AddColumn("NAME", 0.75)
		table.AddColumn("COMMITS", 0.1)
		table.AddColumn("LAST COMMIT", 0.15)

		for _, file := range files {
			table.MustAddRow([]string{
				view.RenderFilename(file.Name),
				strconv.Itoa(file.Count),
				humanize.Time(file.LastCommit)})
		}

		table.Render(v)
	})
}
