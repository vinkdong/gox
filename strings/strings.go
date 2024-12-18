package strings

import "strings"

func Split(s string, seps ...string) []string {
	return genSplit(s, seps, 0, -1)
}
func genSplit(s string, seps []string, start int, n int) []string {
	if len(seps) == 0 || len(s) == 0 {
		return []string{s}
	}

	parts := []string{s}
	for _, sep := range seps {
		var tempParts []string
		for _, part := range parts {
			splitParts := strings.Split(part, sep)
			for _, sp := range splitParts {
				if sp != "" {
					tempParts = append(tempParts, sp)
				}
			}
		}
		parts = tempParts
	}

	return parts
}
