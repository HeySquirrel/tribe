package git

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis/rally"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"strings"
	"time"
)

type Repo struct {
	shell   *shell.Shell
	logger  *tlog.Log
	commits Commits
	Api     *rally.Rally
}

type File struct {
	Name         string
	LastCommit   *Commit
	Contributors Contributors
	Related      []*RelatedFile
	WorkItems    []*rally.Artifact
	Commits      Commits
}

func (f *File) NumberOfDefects() int {
	count := 0
	for _, work := range f.WorkItems {
		if strings.HasPrefix(work.FormattedID, "DE") {
			count += 1
		}
	}
	return count
}

func (f *File) NumberOfStories() int {
	count := 0
	for _, work := range f.WorkItems {
		if strings.HasPrefix(work.FormattedID, "S") {
			count += 1
		}
	}
	return count
}

func (repo *Repo) git(args ...string) (string, error) {
	return repo.shell.Exec("git", args...)
}

func New(dir string, logger *tlog.Log, api *rally.Rally) (*Repo, error) {
	temp := shell.New(dir, logger)
	out, err := temp.Exec("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}

	repo := new(Repo)
	repo.shell = shell.New(strings.TrimSpace(out), logger)
	repo.logger = logger
	repo.Api = api

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	repo.commits, err = CommitsAfter(sixMonthsAgo)
	if err != nil {
		return nil, err
	}

	logger.Add(fmt.Sprintf("Processed %d commits", len(repo.commits)))

	return repo, err
}

func (repo *Repo) Changes() []*File {
	var results = make([]*File, 0)

	cmdOut, err := repo.git("status", "--porcelain")
	if err != nil {
		repo.logger.Add(err.Error())
		return results
	}

	output := strings.Split(cmdOut, "\n")
	for _, change := range output {
		if len(change) > 0 {
			filename := change[3:len(change)]
			results = append(results, repo.GetFile(filename))
		}
	}

	return results
}

func (repo *Repo) GetFile(filename string) *File {
	commits := repo.commits.ContainsFile(filename)
	workItems, _ := repo.Api.GetByFormattedIds(commits.RelatedWorkItems()...)

	file := new(File)
	file.Name = filename
	file.LastCommit = commits[0]
	file.Related = commits.relatedFiles(filename)
	file.Contributors = commits.RelatedContributors()
	file.WorkItems = workItems
	file.Commits = commits

	return file
}
