package git

import (
	"log"
	"os/exec"
	"strings"
)

func Changes() []string {
	var results = make([]string, 1)

	cmdOut, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		log.Panicln(err)
	}

	output := strings.Split(string(cmdOut), "\n")
	for _, change := range output {
		if len(change) > 0 {
			results = append(results, change[3:len(change)])
		}
	}

	return results
}

func FrequentContributors(filename string) []string {
	var results = make([]string, 1)

	command := "git log --pretty=format:'%aN' --follow " + filename + " | sort | uniq -c | sort -rg"
	cmdOut, err := exec.Command("sh", "-c", command).Output()

	if err != nil {
		log.Panicln(err)
	}

	output := strings.Split(string(cmdOut), "\n")
	for _, contributor := range output {
		if len(contributor) > 0 {
			results = append(results, contributor)
		}
	}

	return results
}
