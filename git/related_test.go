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
		commit := new(LogEntry)
		commit.Subject = c.Subject

		entries := Logs{commit}
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
	cases := []struct {
		Author   string
		Expected Contributors
	}{
		{"Bart Simpson", Contributors{NewContributor("Bart Simpson", now)}},
		{"Bart Simpson <bart@simpsons.com>", Contributors{NewContributor("Bart Simpson", now)}},
		{"Bart Simpson and Lisa Simpson",
			Contributors{NewContributor("Bart Simpson", now), NewContributor("Lisa Simpson", now)}},
		{"Homer Simpson, Lisa Simpson and Marge Simpson",
			Contributors{NewContributor("Homer Simpson", now),
				NewContributor("Lisa Simpson", now),
				NewContributor("Marge Simpson", now)}},
	}

	for _, c := range cases {
		commit := new(Commit)
		commit.Author = c.Author
		commit.Date = now

		entries := Commits{commit}
		actual := entries.relatedContributors()

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
