package git

import (
	"regexp"
	"sort"
	"time"
)

type Contributor struct {
	Name       string
	Count      int
	LastCommit *Commit
}

type Contributors []*Contributor

func (c *Contributor) String() string {
	return c.Name
}

func NewContributor(name string, lastCommit *Commit) *Contributor {
	contributor := new(Contributor)
	contributor.Name = name
	contributor.Count = 1
	contributor.LastCommit = lastCommit

	return contributor
}

func (commits *Commits) RelatedWorkItems() []string {
	workItems := make([]string, 0)

	re := regexp.MustCompile("(S|DE|F|s|de|f)[0-9][0-9]+")

	for _, commit := range *commits {
		found := re.FindAllString(commit.Subject, -1)
		if found != nil {
			workItems = append(workItems, found...)
		}
	}

	return workItems
}

func (commits *Commits) RelatedContributors() Contributors {
	contributors := make(Contributors, 0)
	namedContributors := make(map[string]*Contributor)

	remove := regexp.MustCompile(" ?<[^>]+>")
	re := regexp.MustCompile(", | ab?nd |,")

	for _, commit := range *commits {
		authors := remove.ReplaceAllString(commit.Author, "")
		names := re.Split(authors, -1)
		for _, name := range names {
			contributor, ok := namedContributors[name]
			if ok {
				contributor.Count += 1
			} else {
				contributor := NewContributor(name, commit)

				namedContributors[name] = contributor
				contributors = append(contributors, contributor)
			}
		}
	}

	return contributors
}

type RelatedFile struct {
	Name       string
	Count      int
	LastCommit time.Time
}

type byRelevance []*RelatedFile

func (a byRelevance) Len() int      { return len(a) }
func (a byRelevance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRelevance) Less(i, j int) bool {
	return a[i].Count < a[j].Count
}

func NewRelatedFile(name string, lastCommitTime time.Time) *RelatedFile {
	relatedFile := new(RelatedFile)
	relatedFile.Name = name
	relatedFile.Count = 1
	relatedFile.LastCommit = lastCommitTime

	return relatedFile
}

func (commits *Commits) relatedFiles(filename string) []*RelatedFile {
	files := make([]*RelatedFile, 0)
	namedFiles := make(map[string]*RelatedFile)

	for _, commit := range *commits {
		for _, file := range commit.Files {
			if file == filename {
				continue
			}

			relatedFile, ok := namedFiles[file]
			if ok {
				relatedFile.Count += 1
			} else {
				relatedFile := NewRelatedFile(file, commit.Date)

				namedFiles[file] = relatedFile
				files = append(files, relatedFile)
			}
		}
	}

	sort.Sort(sort.Reverse(byRelevance(files)))
	return files
}
