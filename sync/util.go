package sync

import (
	"strconv"
	"strings"
)

func ExtractIdFromUrl(url string) int {
	trimmed := strings.TrimRight(url, "/")
	parts := strings.Split(trimmed, "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])
	return id
}
