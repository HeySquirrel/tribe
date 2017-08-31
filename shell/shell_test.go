package shell

import (
	"strings"
	"testing"
)

func TestRunCommandFailures(t *testing.T) {
	shell := NewInWd()

	cases := []struct {
		Cmd      string
		Expected string
	}{
		{"ls 1231232132", "1231232132"},
		{"cat 12332121", "12332121"},
	}
	for _, c := range cases {
		cmd := c.Cmd
		expected := c.Expected
		_, err := shell.Exec(cmd)
		if !strings.Contains(err.Error(), expected) {
			t.Error("'" + cmd + "' did not error correctly, instead err was " + err.Error())
		}
	}
}

func TestRunCommandSuccess(t *testing.T) {
	shell := NewInWd()

	cases := []struct {
		Cmd string
	}{
		{"ls"},
		{"whoami"},
	}
	for _, c := range cases {
		cmd := c.Cmd
		_, err := shell.Exec(cmd)
		if err != nil {
			t.Error("'" + cmd + "' errored, err was " + err.Error())
		}
	}
}
