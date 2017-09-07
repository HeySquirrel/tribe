package git

import (
	"regexp"
)

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
