package model

import (
	"bufio"
	"fmt"
	"os"
)

type Blame struct {
	File  string
	Start int
	End   int
	Lines []*Line
}

type Line struct {
	Text   string
	Number int
}

func NewBlame(filename string, start, end int) (*Blame, error) {
	if start <= 0 || end <= 0 {
		return nil, fmt.Errorf("fatal: invalid line numbers %d:%d", start, end)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]*Line, 0)
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		lines = append(lines, &Line{Text: scanner.Text(), Number: i})
	}

	numberOfLines := len(lines)
	if numberOfLines < start || numberOfLines < end {
		return nil, fmt.Errorf(
			"fatal: file %s has only %d lines",
			filename,
			numberOfLines,
		)
	}

	return &Blame{File: filename, Start: start, End: end, Lines: lines}, nil
}
