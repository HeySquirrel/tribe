package git

import (
	"fmt"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Contributor struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     int
}

type RelatedFile struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     int
}

type byRelevance []*RelatedFile

func (a byRelevance) Len() int      { return len(a) }
func (a byRelevance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRelevance) Less(i, j int) bool {
	if a[i].UnixTime == a[j].UnixTime {
		return a[i].Count < a[j].Count
	}

	return a[i].UnixTime < a[j].UnixTime
}

type Repo struct {
	shell  *shell.Shell
	logger *tlog.Log
}

func (repo *Repo) git(args ...string) (string, error) {
	return repo.shell.Exec("git", args...)
}

func (repo *Repo) log(args ...string) (string, error) {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	after := fmt.Sprintf("--after=%s", sixMonthsAgo.Format("2006/01/02"))

	logCommand := make([]string, 0)
	logCommand = append(logCommand, "log", "--no-merges", after)
	logCommand = append(logCommand, args...)

	results, err := repo.git(logCommand...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(results), nil
}

func New(dir string, logger *tlog.Log) (*Repo, error) {
	temp := shell.New(dir, logger)
	out, err := temp.Exec("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}

	repo := new(Repo)
	repo.shell = shell.New(strings.TrimSpace(out), logger)
	repo.logger = logger

	return repo, err
}

func (repo *Repo) Changes() []string {
	var results = make([]string, 1)

	cmdOut, err := repo.git("status", "--porcelain")
	if err != nil {
		repo.logger.Add(err.Error())
		return results
	}

	output := strings.Split(cmdOut, "\n")
	for _, change := range output {
		if len(change) > 0 {
			results = append(results, change[3:len(change)])
		}
	}

	return results
}

type LogEntry struct {
	Sha          string
	Subject      string
	Author       string
	RelativeDate string
	UnixTime     int
}

func (repo *Repo) Related(filename string) ([]*RelatedFile, []string, []*Contributor) {
	logEntries := make([]*LogEntry, 0)
	out, err := repo.log("--pretty=format:%H%m%s%m%aN%m%ar%m%at", "--follow", filename)
	if err != nil {
		repo.logger.Add(err.Error())
	}

	logs := strings.Split(out, "\n")
	for _, log := range logs {
		if len(log) == 0 {
			continue
		}

		parts := strings.Split(log, ">")
		entry := new(LogEntry)
		entry.Sha = parts[0]
		entry.Subject = parts[1]
		entry.Author = strings.TrimSpace(parts[2])
		entry.RelativeDate = parts[3]
		entry.UnixTime, err = strconv.Atoi(parts[4])
		if err != nil {
			repo.logger.Add(err.Error())
		}

		logEntries = append(logEntries, entry)
	}

	return repo.relatedFiles(logEntries, filename),
		repo.relatedWorkItems(logEntries),
		repo.relatedContributors(logEntries)
}

func (repo *Repo) relatedFiles(entries []*LogEntry, filename string) []*RelatedFile {
	files := make([]*RelatedFile, 0)
	namedFiles := make(map[string]*RelatedFile)

	for _, entry := range entries {
		out, err := repo.git("show", "--pretty=format:%ar%m%at", "--name-only", entry.Sha)
		if err != nil {
			repo.logger.Add(err.Error())
		}
		lines := strings.Split(out, "\n")
		dateData := strings.Split(lines[0], ">")

		for _, file := range lines[1:] {
			if len(strings.TrimSpace(file)) == 0 || file == filename {
				continue
			}

			relatedFile, ok := namedFiles[file]
			if ok {
				relatedFile.Count += 1
			} else {
				relatedFile := new(RelatedFile)
				relatedFile.Name = file
				relatedFile.Count = 1
				relatedFile.RelativeDate = dateData[0]
				relatedFile.UnixTime, err = strconv.Atoi(dateData[1])
				if err != nil {
					repo.logger.Add(err.Error())
				}

				namedFiles[file] = relatedFile
				files = append(files, relatedFile)
			}
		}
	}

	sort.Sort(sort.Reverse(byRelevance(files)))
	return files
}

func (repo *Repo) relatedWorkItems(entries []*LogEntry) []string {
	workItems := make([]string, 0)

	re := regexp.MustCompile("(S|DE)[0-9][0-9]+")

	for _, entry := range entries {
		found := re.FindString(entry.Subject)
		if len(found) > 0 {
			workItems = append(workItems, found)
		}
	}

	return workItems
}

func (repo *Repo) relatedContributors(entries []*LogEntry) []*Contributor {
	contributors := make([]*Contributor, 0)
	namedContributors := make(map[string]*Contributor)

	for _, entry := range entries {
		name := entry.Author

		contributor, ok := namedContributors[name]
		if ok {
			contributor.Count += 1
		} else {
			contributor := new(Contributor)
			contributor.Name = name
			contributor.Count = 1
			contributor.RelativeDate = entry.RelativeDate
			contributor.UnixTime = entry.UnixTime

			namedContributors[name] = contributor
			contributors = append(contributors, contributor)
		}
	}

	return contributors
}
