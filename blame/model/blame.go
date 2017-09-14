package model

import (
	"bufio"
	"os"
)

type Blame struct {
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

	return &Blame{Lines: lines}, nil
}
