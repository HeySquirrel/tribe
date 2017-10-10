package risk

import (
	"fmt"
	"io"
	"time"

	"github.com/HeySquirrel/tribe/git"
	humanize "github.com/dustin/go-humanize"
)

// Risk is an attempt to quantify the risk of changing a file.
type Risk struct {
	file                     string
	workItemCount            int
	contributorCount         int
	commitCount              int
	commitCountLastSixMonths int
	lastCommit               time.Time
	score                    float64
}

// Calculate the risk of the given file.
func Calculate(file string, commits git.Commits) Risk {
	risk := Risk{}
	risk.file = file

	fileCommits := commits.ContainsFile(file)

	sixMonths := time.Now().AddDate(0, -6, 0)
	risk.commitCountLastSixMonths = fileCommits.CountAfter(sixMonths)
	risk.workItemCount = len(fileCommits.RelatedItems())
	risk.contributorCount = len(fileCommits.RelatedContributors())
	risk.commitCount = len(fileCommits)
	risk.lastCommit = fileCommits[0].Date
	risk.score = (0.33 * risk.workItemCountScore()) +
		(0.33 * risk.contributorCountScore()) +
		(0.33 * risk.commitScore())

	return risk
}

// Write the results of the risk assessment
func (r *Risk) Write(writer io.Writer) {
	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "Risk for '%s'\n\n", r.file)
	fmt.Fprintf(writer, "%20s - Last commit\n", humanize.Time(r.lastCommit))
	fmt.Fprintf(writer, "%20d - Commit count\n", r.commitCount)
	fmt.Fprintf(writer, "%20d - Commits last six months\n", r.commitCountLastSixMonths)
	fmt.Fprintf(writer, "%20d - Work items\n", r.workItemCount)
	fmt.Fprintf(writer, "%20d - Contributors\n", r.contributorCount)
	fmt.Fprintf(writer, "%16s%.2f - Risk\n", "", r.score)
	fmt.Fprintln(writer)

}

func (r *Risk) workItemCountScore() float64 {
	riskyCount := 20
	var divisor int
	if r.workItemCount > riskyCount {
		divisor = r.workItemCount
	} else {
		divisor = riskyCount
	}

	return float64(r.workItemCount) / float64(divisor)
}

func (r *Risk) contributorCountScore() float64 {
	riskyCount := 20
	var divisor int
	if r.contributorCount > riskyCount {
		divisor = r.contributorCount
	} else {
		divisor = riskyCount
	}

	return float64(r.contributorCount) / float64(divisor)
}

func (r *Risk) commitScore() float64 {
	count := r.commitCount - r.commitCountLastSixMonths
	return float64(count) / float64(r.commitCount)
}
