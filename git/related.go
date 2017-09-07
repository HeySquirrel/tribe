package git

import (
	humanize "github.com/dustin/go-humanize"
	"regexp"
	"time"
)

type Contributor struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     time.Time
}

type Contributors []*Contributor

func (c *Contributor) String() string {
	return c.Name
}

func NewContributor(name string, lastContribution time.Time) *Contributor {
	contributor := new(Contributor)
	contributor.Name = name
	contributor.Count = 1
	contributor.UnixTime = lastContribution
	contributor.RelativeDate = humanize.Time(lastContribution)

	return contributor
}

func (entries *Logs) relatedWorkItems() []string {
	workItems := make([]string, 0)

	re := regexp.MustCompile("(S|DE|F)[0-9][0-9]+")

	for _, entry := range *entries {
		found := re.FindAllString(entry.Subject, -1)
		if found != nil {
			workItems = append(workItems, found...)
		}
	}

	return workItems
}

func (entries *Logs) relatedContributors() Contributors {
	contributors := make(Contributors, 0)
	namedContributors := make(map[string]*Contributor)

	remove := regexp.MustCompile(" ?<[^>]+>")
	re := regexp.MustCompile(", | and |,")

	for _, entry := range *entries {
		authors := remove.ReplaceAllString(entry.Author, "")
		names := re.Split(authors, -1)
		for _, name := range names {
			contributor, ok := namedContributors[name]
			if ok {
				contributor.Count += 1
			} else {
				contributor := NewContributor(name, entry.UnixTime)

				namedContributors[name] = contributor
				contributors = append(contributors, contributor)
			}
		}
	}

	return contributors
}
