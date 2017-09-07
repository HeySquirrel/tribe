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

func (entries *Logs) relatedContributors() []*Contributor {
	contributors := make([]*Contributor, 0)
	namedContributors := make(map[string]*Contributor)

	for _, entry := range *entries {
		name := entry.Author

		contributor, ok := namedContributors[name]
		if ok {
			contributor.Count += 1
		} else {
			contributor := new(Contributor)
			contributor.Name = name
			contributor.Count = 1
			contributor.UnixTime = entry.UnixTime
			contributor.RelativeDate = humanize.Time(entry.UnixTime)

			namedContributors[name] = contributor
			contributors = append(contributors, contributor)
		}
	}

	return contributors
}
