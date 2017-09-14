package model

import (
	"bufio"
	"os"
)

type Blame struct {
	File  string
	Start int
	End   int
	Lines []Line
}

type Line struct {
	Text   string
	Number int
}

func New(filename string) (*Blame, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]Line, 20)
	scanner := bufio.NewScanner(file)
	for i := 0; i < 20 && scanner.Scan(); i++ {
		lines[i] = Line{Text: scanner.Text(), Number: i}
	}

	return &Blame{File: filename, Start: 0, End: 20, Lines: lines}, nil
}
