package model

import (
	"bufio"
	"fmt"
	"github.com/heysquirrel/tribe/apis"
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

type History struct {
	commits   git.Commits
	workitems []apis.WorkItem
	Start     int
	End       int
}

func (h *History) GetCommits() git.Commits           { return h.commits }
func (h *History) GetWorkItems() []apis.WorkItem     { return h.workitems }
func (h *History) GetContributors() git.Contributors { return h.commits.RelatedContributors() }

type Line struct {
	Filename string
	Text     string
	Number   int
	once     sync.Once
	history  *History
}

func NewHistory(line *Line) *History {
	start := 1
	end := line.Number + 1

	if line.Number > 1 {
		start = line.Number - 1
	}

	commits, err := git.Log(fmt.Sprintf("-L%d,%d:%s", start, end, line.Filename))
	if err != nil {
		log.Panicln(err)
	}

	return &History{commits: commits, Start: start, End: end}
}

func (l *Line) gatherHistory() {
	l.history = NewHistory(l)
}

func (l *Line) GetHistory() <-chan *History {
	c := make(chan *History)
	go func(line *Line) {
		defer close(c)
		line.once.Do(line.gatherHistory)
		c <- line.history
	}(l)

	return c
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
