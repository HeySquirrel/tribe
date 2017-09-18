package model

import (
	"bufio"
	"fmt"
	"github.com/heysquirrel/tribe/git"
	"os"
)

type File struct {
	Filename string
	Start    int
	End      int
	Lines    []*Line
}

type Line struct {
	Text   string
	Number int
}

func NewFile(filename string, start, end int) (*File, error) {
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

	return &File{Filename: filename, Start: start, End: end, Lines: lines}, nil
}

func (f *File) Len() int {
	return len(f.Lines)
}

func (f *File) GetLine(lineNumber int) *Line {
	return f.Lines[lineNumber-1]
}

func (f *File) Blame(start, end int) (git.Commits, error) {
	return git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, f.Filename))

}
