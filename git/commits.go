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

type Commit struct {
	Sha     string
	Subject string
	Body    string
	Author  string
	Date    time.Time
	Files   []string
}

type Commits []*Commit

func (commit *Commit) HasFile(filename string) bool {
	for _, file := range commit.Files {
		if file == filename {
			return true
		}
	}

	return false
}

func (l *Commits) ContainsFile(filename string) Commits {
	commits := make(Commits, 0)

	for _, commit := range *l {
		if commit.HasFile(filename) {
			commits = append(commits, commit)
		}
	}

	return commits
}

func (repo *Repo) CommitsAfter(after time.Time) (Commits, error) {
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

	commits := parse(stdout)

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func parse(reader io.Reader) Commits {
	commits := make([]*Commit, 0)
	scanner := bufio.NewScanner(reader)

	var currentCommit *Commit
	var body []string

	for scanner.Scan() {
		currentLine := scanner.Text()
		if len(currentLine) == 0 {
			if len(body) != 0 && currentCommit.Body == "" {
				currentCommit.Body = strings.Join(body, " ")
			}
			continue
		}

		if strings.HasPrefix(currentLine, commit) {
			currentCommit = new(Commit)
			body = make([]string, 0)

			commits = append(commits, currentCommit)

			currentCommit.Sha = strings.TrimPrefix(currentLine, commit)
			currentCommit.Files = make([]string, 0)
			continue
		}

		if strings.HasPrefix(currentLine, author) {
			authorLine := strings.TrimPrefix(currentLine, author)
			address, err := mail.ParseAddress(authorLine)
			if err != nil {
				currentCommit.Author = authorLine
			} else {
				currentCommit.Author = address.Name
			}

			continue
		}

		if strings.HasPrefix(currentLine, date) {
			dateLine := strings.TrimPrefix(currentLine, date)
			timestr := strings.TrimSpace(strings.Split(dateLine, "-")[0])

			unixTime, err := strconv.ParseInt(timestr, 10, 64)
			if err != nil {
				currentCommit.Date = time.Unix(0, 0)
			} else {
				currentCommit.Date = time.Unix(unixTime, 0)
			}

			continue
		}

		if strings.HasPrefix(currentLine, "    ") {
			line := strings.TrimSpace(currentLine)
			if currentCommit.Subject == "" {
				currentCommit.Subject = line
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
			currentCommit.Files = append(currentCommit.Files, lineScanner.Text())
			continue
		}

	}

	return commits
}
