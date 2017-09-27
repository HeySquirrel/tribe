package widgets

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/config"
	"github.com/heysquirrel/tribe/git"
	"github.com/heysquirrel/tribe/work"
	"io"
	"regexp"
)

type WorkItems []work.Item

func (items WorkItems) Display(writer io.Writer) {
	for _, item := range items {
		fmt.Fprintf(writer, "%10s - %s\n", item.GetId(), item.GetName())
	}
}
func (items WorkItems) Len() int { return len(items) }

type ContributorItems []*git.Contributor

func (c ContributorItems) Display(writer io.Writer) {
	for _, contributor := range c {
		fmt.Fprintf(writer, "  %-20s - %d Commits - %s\n",
			contributor.Name,
			contributor.Count,
			humanize.Time(contributor.LastCommit.Date),
		)
	}
}
func (c ContributorItems) Len() int { return len(c) }

type CommitItems []*git.Commit

func (c CommitItems) Display(writer io.Writer) {
	revert := regexp.MustCompile("(r|R)evert")
	magenta := func(s string) string { return color.MagentaString(s) }
	cyan := func(s string) string { return color.CyanString(s) }

	for _, commit := range c {
		subject := commit.Subject
		for _, re := range config.Matchers() {
			subject = re.ReplaceAllStringFunc(subject, magenta)
		}
		subject = revert.ReplaceAllStringFunc(subject, cyan)

		fmt.Fprintf(writer, " %10s - %s - %s\n",
			commit.Sha[0:9],
			subject,
			humanize.Time(commit.Date),
		)
	}
}
func (c CommitItems) Len() int { return len(c) }

type FileItems model.File

func (f FileItems) Display(writer io.Writer) {
	for _, line := range f.Lines {
		fmt.Fprintf(writer, "%5d| %s\n", line.Number, line.Text)
	}
}
func (f FileItems) Len() int { file := model.File(f); return file.Len() }
