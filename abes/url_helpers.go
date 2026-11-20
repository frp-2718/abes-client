package abes

import "strings"

// concatValues groups values into comma-separated batches. A maxValues below 1
// uses the service limit; a larger limit is capped at the service limit.
func concatValues(values []string, maxValues, serviceMax int) []string {
	if maxValues < 1 {
		maxValues = serviceMax
	} else if maxValues > serviceMax {
		maxValues = serviceMax
	}

	result := []string{}
	for len(values) > maxValues {
		result = append(result, strings.Join(values[:maxValues], ","))
		values = values[maxValues:]
	}
	if len(values) > 0 {
		result = append(result, strings.Join(values, ","))
	}
	return result
}

// joinURLPath joins a base URL and path with exactly one slash at the boundary.
func joinURLPath(base, path string) string {
	if strings.HasSuffix(base, "/") && strings.HasPrefix(path, "/") {
		return base + path[1:]
	}
	if !strings.HasSuffix(base, "/") && !strings.HasPrefix(path, "/") {
		return base + "/" + path
	}
	return base + path
}
