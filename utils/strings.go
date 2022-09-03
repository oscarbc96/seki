package utils

import "strings"

func FindStartAndEndColumn(str, substr string) (int, int) {
	idx := strings.Index(str, substr)
	if idx == -1 {
		return -1, -1
	}
	return idx + 1, idx + len(substr)
}
