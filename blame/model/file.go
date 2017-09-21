package model

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	RelPath string
	Name    string
	Start   int
	End     int
	Lines   []*Line
}

type Line struct {
	File   *File
	Text   string
	Number int
}

func NewFile(filename string, start, end int) (*File, error) {
	if start <= 0 || end <= 0 {
		return nil, fmt.Errorf("fatal: invalid line numbers %d:%d", start, end)
	}

	_, name := filepath.Split(filename)
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	file := new(File)
	file.RelPath = filename
	file.Name = name
	file.Start = start
	file.End = end
	file.Lines = make([]*Line, 0)

	scanner := bufio.NewScanner(reader)
	for i := 1; scanner.Scan(); i++ {
		file.Lines = append(
			file.Lines,
			&Line{File: file, Text: scanner.Text(), Number: i},
		)
	}

	numberOfLines := len(file.Lines)
	if numberOfLines < start || numberOfLines < end {
		return nil, fmt.Errorf(
			"fatal: file %s has only %d lines",
			filename,
			numberOfLines,
		)
	}

	return file, nil
}

func (f *File) Len() int {
	return len(f.Lines)
}

func (f *File) GetLine(lineNumber int) *Line {
	return f.Lines[lineNumber-1]
}
