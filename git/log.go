package git

import (
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	commit = "commit "
	author = "Author: "
	date   = "Date: "
)

type LogEntry struct {
	Sha          string
	Subject      string
	Body         string
	Author       string
	RelativeDate string
	UnixTime     int64
	Files        []string
}

type Logs []*LogEntry

func (entry *LogEntry) HasFile(filename string) bool {
	for _, file := range entry.Files {
		if file == filename {
			return true
		}
	}

	return false
}

func (l *Logs) ContainsFile(filename string) Logs {
	logs := make(Logs, 0)

	for _, entry := range *l {
		if entry.HasFile(filename) {
			logs = append(logs, entry)
		}
	}

	return logs
}

func (repo *Repo) LogsAfter(after time.Time) (Logs, error) {
	afterArg := fmt.Sprintf("--after=%s", after.Format("2006/01/02"))
	args := []string{"log", "--no-merges", "--raw", "--date=unix", afterArg}

	repo.logger.Add("git " + strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	entries := parse(stdout)

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func parse(reader io.Reader) Logs {
	entries := make([]*LogEntry, 0)
	scanner := bufio.NewScanner(reader)

	var currentEntry *LogEntry
	var body []string

	for scanner.Scan() {
		currentLine := scanner.Text()
		if len(currentLine) == 0 {
			if len(body) != 0 && currentEntry.Body == "" {
				currentEntry.Body = strings.Join(body, " ")
			}
			continue
		}

		if strings.HasPrefix(currentLine, commit) {
			currentEntry = new(LogEntry)
			body = make([]string, 0)

			entries = append(entries, currentEntry)

			currentEntry.Sha = strings.TrimPrefix(currentLine, commit)
			currentEntry.Files = make([]string, 0)
			continue
		}

		if strings.HasPrefix(currentLine, author) {
			authorLine := strings.TrimPrefix(currentLine, author)
			address, err := mail.ParseAddress(authorLine)
			if err != nil {
				currentEntry.Author = authorLine
			} else {
				currentEntry.Author = address.Name
			}

			continue
		}

		if strings.HasPrefix(currentLine, date) {
			dateLine := strings.TrimPrefix(currentLine, date)
			time := strings.TrimSpace(strings.Split(dateLine, "-")[0])

			unixTime, err := strconv.ParseInt(time, 10, 64)
			if err != nil {
				currentEntry.UnixTime = 0
			} else {
				currentEntry.UnixTime = unixTime
			}

			continue
		}

		if strings.HasPrefix(currentLine, "    ") {
			line := strings.TrimSpace(currentLine)
			if currentEntry.Subject == "" {
				currentEntry.Subject = line
			} else {
				body = append(body, line)
			}

			continue
		}

		if strings.HasPrefix(currentLine, ":") {
			lineReader := strings.NewReader(strings.TrimPrefix(currentLine, ":"))
			lineScanner := bufio.NewScanner(lineReader)
			lineScanner.Split(bufio.ScanWords)

			lineScanner.Scan() // Old file permissions
			lineScanner.Scan() // New file permissions
			lineScanner.Scan() // Old Tree
			lineScanner.Scan() // New Tree
			lineScanner.Scan() // Change Types

			lineScanner.Scan() // File Name
			currentEntry.Files = append(currentEntry.Files, lineScanner.Text())
			continue
		}

	}

	return entries
}
