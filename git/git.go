package git

import (
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"strings"
	"time"
)

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
