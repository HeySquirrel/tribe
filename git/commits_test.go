package git

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOneRawGitLog(t *testing.T) {
	input, cleanup := setup(t, "single_log_entry.txt")
	defer cleanup()

	entries := parse(input)

	AssertInt(t, len(entries), 1)
	AssertString(t, entries[0].Sha, "2366697ee5e9c56106a595fdc84c4218f7f9db04")
	AssertString(t, entries[0].Author, "Joe Developer")
	AssertInt64(t, entries[0].Date.Unix(), 1502491416)
	AssertString(t, entries[0].Subject, "DE3456 - All logins should respect login bypass_login")
	AssertString(t, entries[0].Body, "/oauth/token wasn't respecting the organization.bypass_login flag. Unlike the web login /oauth/token was being controlled by User.authenticate. We changed User.authenticate to respect the bypass_login flag.")
	AssertString(t, entries[0].Files[0], "app/models/user.rb")
}

func TestManyRawGitLog(t *testing.T) {
	input, cleanup := setup(t, "multiple_files_per_commit_log_entries.txt")
	defer cleanup()

	entries := parse(input)

	AssertInt(t, len(entries), 8)
	AssertString(t, entries[3].Sha, "306f80d5edba8d9e6dd391916604ff9441898dab")
	AssertInt(t, len(entries[3].Files), 7)
}

func TestDifferentGitLogs(t *testing.T) {
	input, cleanup := setup(t, "different_log_entries.txt")
	defer cleanup()

	entries := parse(input)

	AssertInt(t, len(entries), 2)
	AssertString(t, "log/log.go", entries[0].Files[0])
}

func TestHasFile(t *testing.T) {
	input, cleanup := setup(t, "single_log_entry.txt")
	defer cleanup()

	entries := parse(input)

	AssertInt(t, len(entries), 1)
	if entries[0].HasFile("app/models/user.rb") != true {
		t.Error("Expected to fine user.rb")
	}

	if entries[0].HasFile("unknown/unknown.rb") {
		t.Error("Unknown file found")
	}
}

func TestContainsFile(t *testing.T) {
	input, cleanup := setup(t, "multiple_files_per_commit_log_entries.txt")
	defer cleanup()

	entries := parse(input)

	coolEntries := entries.ContainsFile("app/models/cool/cool.rb")
	AssertInt(t, len(coolEntries), 1)

	unknownEntries := entries.ContainsFile("unkown/unknown.rb")
	AssertInt(t, len(unknownEntries), 0)
}

func TestLogContainsDiff(t *testing.T) {
	input, cleanup := setup(t, "expanded_log_entries.txt")
	defer cleanup()

	entries := parse(input)

	AssertInt(t, len(entries), 2)
	AssertString(t, "1633ea32bf23ba1f59626d4f2330cf65c4a2ec8d", entries[0].Sha)
}

func TestLogsAfter(t *testing.T) {
	logs, err := CommitsAfter(time.Now().AddDate(0, -6, 0))
	if err != nil {
		t.Fatalf("unable to get logs: %v", err)
	}

	if len(logs) == 0 {
		t.Error("Expected more than 0 logs")
	}
}

func setup(t *testing.T, fileFixture string) (io.Reader, func()) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to get pwd")
	}

	fixture := filepath.Join(pwd, "fixtures", fileFixture)
	log, err := os.Open(fixture)
	if err != nil {
		t.Fatal("Unable to read fixture")
	}

	return log, func() { log.Close() }
}

func AssertInt(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Errorf("Expected: %d, But got: %d", expected, actual)
	}
}

func AssertInt64(t *testing.T, actual int64, expected int64) {
	if actual != expected {
		t.Errorf("Expected: %d, But got: %d", expected, actual)
	}
}

func AssertString(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("Expected: '%s', But got: '%s'", expected, actual)
	}
}
