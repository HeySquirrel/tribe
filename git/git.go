package git

import (
	humanize "github.com/dustin/go-humanize"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"sort"
	"strings"
	"time"
)

type RelatedFile struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     time.Time
}

type byRelevance []*RelatedFile

func (a byRelevance) Len() int      { return len(a) }
func (a byRelevance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRelevance) Less(i, j int) bool {
	if a[i].UnixTime == a[j].UnixTime {
		return a[i].Count < a[j].Count
	}

	return a[i].UnixTime.Before(a[j].UnixTime)
}

type Repo struct {
	shell  *shell.Shell
	logger *tlog.Log
}

func (repo *Repo) git(args ...string) (string, error) {
	return repo.shell.Exec("git", args...)
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
	var results = make([]string, 0)

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

func (repo *Repo) Related(filename string) ([]*RelatedFile, []string, []*Contributor) {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	logs, err := repo.LogsAfter(sixMonthsAgo)
	if err != nil {
		repo.logger.Add(err.Error())
	}
	relatedLogs := logs.ContainsFile(filename)

	return relatedLogs.relatedFiles(filename), relatedLogs.relatedWorkItems(), relatedLogs.relatedContributors()
}

func (entries *Logs) relatedFiles(filename string) []*RelatedFile {
	files := make([]*RelatedFile, 0)
	namedFiles := make(map[string]*RelatedFile)

	for _, entry := range *entries {
		for _, file := range entry.Files {
			if file == filename {
				continue
			}

			relatedFile, ok := namedFiles[file]
			if ok {
				relatedFile.Count += 1
			} else {
				relatedFile := new(RelatedFile)
				relatedFile.Name = file
				relatedFile.Count = 1
				relatedFile.UnixTime = entry.UnixTime
				relatedFile.RelativeDate = humanize.Time(entry.UnixTime)

				namedFiles[file] = relatedFile
				files = append(files, relatedFile)
			}
		}
	}

	sort.Sort(sort.Reverse(byRelevance(files)))
	return files
}
