package sync

import (
	"sort"
	"strconv"
	"strings"
)

type HasId interface {
	GetId() int
}

func ExtractIdFromUrl(url string) int {
	trimmed := strings.TrimRight(url, "/")
	parts := strings.Split(trimmed, "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])
	return id
}

func SortById[T HasId](items []T) []T {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetId() < items[j].GetId()
	})
	return items
}
