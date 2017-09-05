package git

import (
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func (repo *Repo) RelatedFiles(filename string) []*RelatedFile {
	files := make([]*RelatedFile, 0)
	namedFiles := make(map[string]*RelatedFile)

	if len(filename) == 0 {
		return files
	}

	out, err := repo.git("log", "--pretty=format:%H", "--follow", filename)
	if err != nil {
		repo.logger.Add(err.Error())
	}

	output := strings.Split(out, "\n")
	for _, sha := range output {
		out, err = repo.git("show", "--pretty=format:%ar%m%at", "--name-only", sha)
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

func (repo *Repo) RelevantWorkItems(filename string) []string {
	workItems := make([]string, 0)

	if len(filename) == 0 {
		return workItems
	}

	out, err := repo.git("log", "--pretty=format:%s", "--follow", filename)
	if err != nil {
		repo.logger.Add(err.Error())
	}

	subjects := strings.Split(out, "\n")
	re := regexp.MustCompile("(S|DE)[0-9][0-9]+")

	for _, subjects := range subjects {
		found := re.FindString(subjects)
		if len(found) > 0 {
			workItems = append(workItems, found)
		}
	}

	return workItems
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
