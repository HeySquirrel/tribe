package view

import (
	"testing"
)

func TestRenderFilename(t *testing.T) {
	cases := []struct {
		filename string
		expected string
	}{
		{"foo/bar/baz/model/user.rb", ".../model/user.rb"},
		{"bart.js", "bart.js"},
		{"model/user.rb", "model/user.rb"},
		{"app/model/user.rb", ".../model/user.rb"},
	}

	for _, c := range cases {
		if c.expected != RenderFilename(c.filename) {
			t.Errorf("Expected: '%s', Got: '%s'", c.expected, RenderFilename(c.filename))
		}
	}
}
