package app

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/view"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
	"strings"
)

type Column struct {
	name string
	size float64
}

type Row []string

type Table struct {
	width   int
	columns []Column
	rows    []Row
}

func NewColumn(name string, size float64) Column {
	return Column{name: name, size: size}
}

func NewTable(width int, columns ...Column) *Table {
	table := new(Table)
	table.width = width
	table.columns = columns
	table.rows = make([]Row, 0)

	return table
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
		columnSize := int(float64(maxView) * column.size)
		header = append(header, Center2(column.name, columnSize))
	}
	fmt.Fprintln(w, strings.Join(header, "|"))

	fmt.Fprintf(w, "+%s+\n", strings.Repeat("-", maxView))

	for _, row := range t.rows {
		columns := make([]string, 0)

		for j, column := range t.columns {
			columnSize := int(float64(maxView) * column.size)
			data := row[j]
			if j == 0 {
				data = fmt.Sprintf(" 🌶  %s", data)
			}
			columnFormat := fmt.Sprintf(" %%-%ds", columnSize-1)
			columns = append(columns, fmt.Sprintf(columnFormat, data))
		}

		fmt.Fprintln(w, strings.Join(columns, "|"))
	}
}

func (a *App) UpdateContributors2(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		maxX, _ := v.Size()
		table := NewTable(maxX, NewColumn("NAME", 0.55), NewColumn("COMMITS", 0.2), NewColumn("LAST COMMIT", 0.25))

		for _, contributor := range contributors {

			table.MustAddRow([]string{contributor.Name, strconv.Itoa(contributor.Count), humanize.Time(contributor.LastCommit)})
		}

		table.Render(v)
	})
}

func (a *App) UpdateContributors(contributors []*git.Contributor) {
	a.updateView(contributorsView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := tablewriter.NewWriter(v)
		table.SetColWidth(maxX)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, contributor := range contributors {
			table.Append([]string{
				contributor.Name,
				strconv.Itoa(contributor.Count),
				humanize.Time(contributor.LastCommit),
			})
		}

		table.Render()
	})
}

func (a *App) UpdateRelatedFiles(files []*git.RelatedFile) {
	a.updateView(associatedFilesView, func(v *gocui.View) {
		maxX, _ := v.Size()

		table := tablewriter.NewWriter(v)
		table.SetColWidth(maxX)
		table.SetHeader([]string{"Name", "Commits", "Last Commit"})
		table.SetBorder(false)

		for _, file := range files {
			table.Append([]string{
				view.RenderFilename(file.Name),
				strconv.Itoa(file.Count),
				humanize.Time(file.LastCommit),
			})
		}

		table.Render()
	})
}
