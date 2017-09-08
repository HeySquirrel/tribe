package view

import (
	"path/filepath"
	"strings"
)

func RenderFilename(filename string) string {
	parts := strings.Split(filepath.ToSlash(filename), "/")
	if len(parts) < 3 {
		return filename
	}

	importantParts := append([]string{"..."}, parts[len(parts)-2:]...)

	return filepath.Join(importantParts...)
}
