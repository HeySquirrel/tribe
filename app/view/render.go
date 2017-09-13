package view

import (
	"path/filepath"
	"strings"
)

func RenderFilename(width int, filename string) string {
	if len([]rune(filename)) <= width {
		return filename
	}

	actualWidth := width - 3
	parts := strings.Split(filepath.ToSlash(filename), "/")
	reverse(parts)
	importantParts := make([]string, 0)

	currentWidth := 0
	for _, part := range parts {
		possibleNewWidth := currentWidth + len([]rune(part)) + 1

		if possibleNewWidth < actualWidth {
			importantParts = append(importantParts, part)
			currentWidth = possibleNewWidth
		} else {
			break
		}
	}

	reverse(importantParts)

	return filepath.Join(append([]string{"..."}, importantParts...)...)
}

func reverse(list []string) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}
