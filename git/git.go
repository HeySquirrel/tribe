package git

import (
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"log"
	"strconv"
	"strings"
)

type Contributor struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     int
}

type Repo struct {
	shell  *shell.Shell
	logger *tlog.Log
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

func (repo *Repo) git(args ...string) (string, error) {
	return repo.shell.Exec("git", args...)
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

func (repo *Repo) RelatedFiles(filename string) []string {
	files := make([]string, 0)
	out, err := repo.git("log", "--pretty=format:%H", "--follow", filename)
	if err != nil {
		repo.logger.Add(err.Error())
	}

	output := strings.Split(out, "\n")
	for _, sha := range output {
		out, err = repo.git("show", "--pretty=format:", "--name-only", sha)
		files = append(files, strings.Split(out, "\n")...)
	}

	return files
}

func (repo *Repo) RecentContributors(filename string) []*Contributor {
	contributors := make([]*Contributor, 0)
	namedContributors := make(map[string]*Contributor)

	if len(filename) == 0 {
		return contributors
	}

	out, err := repo.git("log", "--pretty=format:%aN%m%ar%m%at", "--follow", filename)
	if err != nil {
		repo.logger.Add(err.Error())
	}

	output := strings.Split(out, "\n")
	for _, line := range output {
		if len(line) > 0 {
			contributorData := strings.Split(line, ">")
			name := strings.TrimSpace(contributorData[0])

			contributor, ok := namedContributors[name]
			if ok {
				contributor.Count += 1
			} else {
				contributor := new(Contributor)
				contributor.Name = name
				contributor.Count = 1
				contributor.RelativeDate = contributorData[1]
				contributor.UnixTime, err = strconv.Atoi(contributorData[2])
				if err != nil {
					log.Panicln(err)
				}

				namedContributors[name] = contributor
				contributors = append(contributors, contributor)
			}
		}
	}

	return contributors
}
