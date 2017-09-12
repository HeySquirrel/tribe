package view

import (
	"fmt"
	"io"
	"strings"
)

type Justification int

const (
	LEFT Justification = 1 + iota
	RIGHT
)

const (
	HOT = `🌶`
)

type Column struct {
	name          string
	size          int
	justification Justification
}

type Row []string

type Table struct {
	width   int
	columns []Column
	rows    []Row
}

func NewColumn(name string, size int, justification Justification) Column {
	return Column{name: name, size: size, justification: justification}
}

func (c *Column) Render(data string) string {
	format := " %%-%ds"
	if c.justification == RIGHT {
		format = "%%%ds "
	}
	return fmt.Sprintf(fmt.Sprintf(format, c.size-1), data)
}

func NewTable(width int) *Table {
	table := new(Table)
	table.width = width
	table.columns = make([]Column, 0)
	table.rows = make([]Row, 0)

	return table
}

func (t *Table) AddColumn(name string, size float64, justification Justification) {
	columnSize := int(float64(t.width) * size)
	t.columns = append(t.columns, NewColumn(name, columnSize, justification))
}

func (t *Table) MustAddRow(row Row) {
	if len(t.columns) != len(row) {
		panic("Row size should match column size")
	}

	t.rows = append(t.rows, row)
}

func center(s string, width int) string {
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
		header = append(header, center(column.name, column.size))
	}
	fmt.Fprintln(w, strings.Join(header, "|"))

	fmt.Fprintf(w, "+%s+\n", strings.Repeat("-", maxView))

	for i, row := range t.rows {
		columns := make([]string, 0)

		for j, column := range t.columns {
			data := row[j]
			if i == 0 && j == 0 {
				data = fmt.Sprintf(" %s  %s", HOT, data)
			}

			columns = append(columns, column.Render(data))
		}

		fmt.Fprintln(w, strings.Join(columns, "|"))
	}
}
