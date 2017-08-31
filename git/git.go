package git

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Contributor struct {
	Name         string
	Count        int
	RelativeDate string
	UnixTime     int
}

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

func RecentContributors(filename string) []*Contributor {
	contributors := make([]*Contributor, 0)
	namedContributors := make(map[string]*Contributor)

	if len(filename) == 0 {
		return contributors
	}

	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)

	command := exec.Command("git", "log", "--pretty=format:%aN%m%ar%m%at", "--follow", filename)
	command.Stdout = &out
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return contributors
	}

	output := strings.Split(out.String(), "\n")
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
