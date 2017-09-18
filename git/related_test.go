package git

import (
	"testing"
	"time"
)

func TestRelatedWorkItems(t *testing.T) {
	cases := []struct {
		Subject  string
		Expected []string
	}{
		{"F198234_team_coolness- scope coolness by user", []string{"F198234"}},
		{"S9028: Make something cool", []string{"S9028"}},
		{"DE9283: user is uncool", []string{"DE9283"}},
		{"S28973: Remove F2938_uncool_users", []string{"S28973", "F2938"}},
		{"No related work", []string{}},
	}

	for _, c := range cases {
		commit := new(Commit)
		commit.Subject = c.Subject

		entries := Commits{commit}
		actual := entries.relatedWorkItems()

		if len(actual) != len(c.Expected) {
			t.Fatalf("'%v' not equal '%v'", actual, c.Expected)
		}

		for i := range actual {
			if actual[i] != c.Expected[i] {
				t.Errorf("'%s' not equal '%s'", actual[i], c.Expected[i])
			}
		}
	}
}

func TestRelatedContributors(t *testing.T) {
	now := time.Now()
	lastCommit := &Commit{Author: "who cares", Date: now}
	cases := []struct {
		Author   string
		Expected Contributors
	}{
		{"Bart Simpson", Contributors{NewContributor("Bart Simpson", lastCommit)}},
		{"Bart Simpson <bart@simpsons.com>", Contributors{NewContributor("Bart Simpson", lastCommit)}},
		{"Bart Simpson and Lisa Simpson",
			Contributors{NewContributor("Bart Simpson", lastCommit), NewContributor("Lisa Simpson", lastCommit)}},
		{"Homer Simpson, Lisa Simpson and Marge Simpson",
			Contributors{NewContributor("Homer Simpson", lastCommit),
				NewContributor("Lisa Simpson", lastCommit),
				NewContributor("Marge Simpson", lastCommit)}},
	}

	for _, c := range cases {
		commit := new(Commit)
		commit.Author = c.Author
		commit.Date = now

		entries := Commits{commit}
		actual := entries.RelatedContributors()

		if len(actual) != len(c.Expected) {
			t.Fatalf("'%v' not equal '%v'", actual, c.Expected)
		}

		for i := range actual {
			if actual[i].Name != c.Expected[i].Name {
				t.Errorf("'%s' not equal '%s'", actual[i].Name, c.Expected[i].Name)
			}
		}
	}
}
