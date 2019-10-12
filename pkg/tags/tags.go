package tags

import (
	"sort"
	"strings"
)

func CombineTags(sourceTags []string, paramsTags []string) []string {

	uniqueTags := make(map[string]bool, 0)

	for _, tag := range sourceTags {
		uniqueTags[tag] = true
	}

	for _, tag := range paramsTags {
		uniqueTags[tag] = true
	}

	tags := make([]string, 0)

	for tag := range uniqueTags {
		tags = append(tags, tag)
	}

	sort.Strings(tags)

	return tags
}

func FormatTags(tags []string) string {
	sortedTags := CombineTags([]string{}, tags)

	return strings.Join(sortedTags, ", ")
}
