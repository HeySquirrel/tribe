package model

import (
	"bufio"
	"fmt"
	"github.com/heysquirrel/tribe/git"
	"log"
	"os"
	"sync"
	"time"
)

type File struct {
	Filename string
	Start    int
	End      int
	Lines    []*Line
}

type Line struct {
	Filename string
	Text     string
	Number   int
	once     sync.Once
	commits  git.Commits
}

func (l *Line) init() {
	start := 1
	end := l.Number + 1

	if l.Number > 1 {
		start = l.Number - 1
	}

	commits, err := git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, l.Filename))
	if err != nil {
		log.Panicln(err)
	}

	l.commits = commits
}

func (l *Line) GetCommits() git.Commits {
	l.once.Do(l.init)
	return l.commits
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
		lines = append(lines, &Line{Filename: filename, Text: scanner.Text(), Number: i})
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

func (f *File) Logs() (git.Commits, error) {
	commits, err := git.CommitsAfter(time.Now().AddDate(-1, 0, 0))
	if err != nil {
		return nil, err
	}

	return commits.ContainsFile(f.Filename), nil
}
