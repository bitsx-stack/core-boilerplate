package utils

import "strings"

func ParseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}