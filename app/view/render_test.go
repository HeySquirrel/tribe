package view

import (
	"testing"
)

func TestRenderFilename(t *testing.T) {
	cases := []struct {
		filename string
		width    int
		expected string
	}{
		{"foo/bar/baz/model/user.rb", 22, ".../baz/model/user.rb"},
		{"foo/bar/baz/model/user.rb", 20, ".../model/user.rb"},
		{"bart.js", 15, "bart.js"},
		{"model/user.rb", 15, "model/user.rb"},
		{"app/model/user.rb", 20, "app/model/user.rb"},
	}

	for _, c := range cases {
		result := RenderFilename(c.width, c.filename)
		if c.expected != result {
			t.Errorf("Expected: '%s', Got: '%s'", c.expected, result)
		}
	}
}
