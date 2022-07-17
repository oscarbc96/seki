package utils

import "strings"

func FindStartAndEndColumn(s, substr string) (int, int) {
	idx := strings.Index(s, substr)
	if idx == -1 {
		return -1, -1
	}
	return idx + 1, idx + 1 + len(substr)
}
